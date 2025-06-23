package api

import (
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/config"
	commonsHttp "github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/abhirup.dandapat/oms/internal/models"
)

type CreateOrderRequest struct {
	TenantID string `json:"tenant_id" binding:"required"`
	SellerID string `json:"seller_id" binding:"required"`
	HubID    string `json:"hub_id"    binding:"required"`
	SKUID    string `json:"sku_id"    binding:"required"`
	Quantity int64  `json:"quantity"  binding:"required,gt=0"`
}

func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	ctx := c.Request.Context()

	// 1) Check inventory in IMS
	baseURL := config.GetString(ctx, "ims.baseUrl")
	url := fmt.Sprintf("%s/inventory?hub_id=%s&sku_ids=%s", baseURL, req.HubID, req.SKUID)

	transport := &stdhttp.Transport{}
	httpClient, err := commonsHttp.NewHTTPClient("order-service", "", transport)
	if err != nil {
		log.DefaultLogger().Errorf("CreateOrder: NewHTTPClient failed: %v", err)
		c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}

	var invs []models.Inventory
	getReq := &commonsHttp.Request{
		Url:     url,
		Timeout: 5 * time.Second,
	}
	if _, err := httpClient.Get(getReq, &invs); err != nil || len(invs) != 1 {
		log.DefaultLogger().Errorf("CreateOrder: IMS GET failed: %v", err)
		c.JSON(stdhttp.StatusServiceUnavailable, gin.H{"error": i18n.Translate(c, "error.inventory_unavailable")})
		return
	}
	inv := invs[0]
	if inv.QuantityOnHand < req.Quantity {
		c.JSON(stdhttp.StatusConflict, gin.H{"error": i18n.Translate(c, "error.insufficient_inventory")})
		return
	}

	// 2) Reserve (decrement) inventory back to IMS
	putReq := &commonsHttp.Request{
		Url: baseURL + "/inventory",
		Body: map[string]interface{}{
			"tenant_id": req.TenantID,
			"hub_id":    req.HubID,
			"sku_id":    req.SKUID,
			"quantity":  inv.QuantityOnHand - req.Quantity,
		},
		Timeout: 5 * time.Second,
	}
	if _, err := httpClient.Post(putReq, nil); err != nil {
		log.DefaultLogger().Errorf("CreateOrder: IMS PUT failed: %v", err)
		c.JSON(stdhttp.StatusServiceUnavailable, gin.H{"error": i18n.Translate(c, "error.inventory_update_failed")})
		return
	}

	mongoURI := config.GetString(ctx, "mongo.uri")
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.DefaultLogger().Errorf("CreateOrder: mongo connect: %v", err)
		c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	defer cli.Disconnect(ctx)
	coll := cli.Database("omsdb").Collection("orders")

	orderID := uuid.New().String()
	now := time.Now().UTC()
	order := models.Order{
		ID:        orderID,
		TenantID:  req.TenantID,
		SellerID:  req.SellerID,
		HubID:     req.HubID,
		SKUID:     req.SKUID,
		Quantity:  req.Quantity,
		Status:    "new_order",
		CreatedAt: now,
	}
	if _, err := coll.InsertOne(ctx, order); err != nil {
		log.DefaultLogger().Errorf("CreateOrder: mongo insert: %v", err)
		c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}

	producer := kafka.NewProducer(
		kafka.WithBrokers(config.GetStringSlice(ctx, "kafka.brokers")),
		kafka.WithClientID(config.GetString(ctx, "kafka.clientId")+"-producer"),
		kafka.WithKafkaVersion(config.GetString(ctx, "kafka.version")),
	)
	evt := models.OrderCreated{
		OrderID:   orderID,
		TenantID:  req.TenantID,
		SellerID:  req.SellerID,
		HubID:     req.HubID,
		SKUID:     req.SKUID,
		Quantity:  req.Quantity,
		CreatedAt: now,
	}
	payload, _ := json.Marshal(evt)
	msg := &pubsub.Message{
		Topic: config.GetString(ctx, "kafka.topicOrderCreated"),
		Key:   orderID,
		Value: payload,
	}
	if err := producer.Publish(ctx, msg); err != nil {
		log.DefaultLogger().Errorf("CreateOrder: publish order.created failed: %v", err)
	}

	c.JSON(stdhttp.StatusCreated, gin.H{"order_id": orderID})
}
