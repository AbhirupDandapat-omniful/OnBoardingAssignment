package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/abhirup.dandapat/oms/internal/models"
)

func ListOrders(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID := c.Query("tenant_id")
	sellerID := c.Query("seller_id")
	status := c.Query("status")
	from := c.Query("from")
	to := c.Query("to")

	mongoURI := config.GetString(ctx, "mongo.uri")
	clientOpts := options.Client().ApplyURI(mongoURI)
	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.DefaultLogger().Errorf("ListOrders: mongo connect error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	defer mongoClient.Disconnect(ctx)

	coll := mongoClient.Database("omsdb").Collection("orders")

	filter := bson.M{}
	if tenantID != "" {
		filter["tenant_id"] = tenantID
	}
	if sellerID != "" {
		filter["seller_id"] = sellerID
	}
	if status != "" {
		filter["status"] = status
	}
	if from != "" || to != "" {
		dateFilter := bson.M{}
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			dateFilter["$gte"] = t
		}
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			dateFilter["$lte"] = t
		}
		filter["created_at"] = dateFilter
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.DefaultLogger().Errorf("ListOrders: find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	for cursor.Next(ctx) {
		var o models.Order
		if err := cursor.Decode(&o); err != nil {
			log.DefaultLogger().Errorf("ListOrders: decode error: %v", err)
			continue
		}
		orders = append(orders, o)
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
