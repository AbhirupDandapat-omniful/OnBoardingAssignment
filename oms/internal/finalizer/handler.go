package finalizer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	commonsHttp "github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/abhirup.dandapat/oms/internal/models"
)

const updatedTopic = "order.updated"

// Handler finalizes an order: checks/deducts IMS, updates Mongo, emits order.updated.
type Handler struct {
	coll      *mongo.Collection
	client    *commonsHttp.Client
	publisher pubsub.Publisher
	logger    *log.Logger
}

// NewHandler constructs a Handler.
func NewHandler(
	coll *mongo.Collection,
	client *commonsHttp.Client,
	publisher pubsub.Publisher,
) *Handler {
	return &Handler{
		coll:      coll,
		client:    client,
		publisher: publisher,
		logger:    log.DefaultLogger(),
	}
}

// Process implements pubsub.IPubSubMessageHandler.
func (h *Handler) Process(ctx context.Context, msg *pubsub.Message) error {
	// 1) Decode order.created payload
	var oc models.OrderCreated
	if err := json.Unmarshal(msg.Value, &oc); err != nil {
		h.logger.Errorf("invalid payload: %v", err)
		return err
	}
	h.logger.Infof("Finalizing order %s", oc.OrderID)

	// 2) GET inventory from IMS
	getReq := &commonsHttp.Request{
		Url:     fmt.Sprintf("/inventory?hub_id=%s&sku_ids=%s", oc.HubID, oc.SKUID),
		Timeout: 5 * time.Second,
	}
	getResp, err := h.client.Get(getReq, nil)
	if err != nil {
		h.logger.Errorf("IMS GET error: %v", err)
		return err
	}
	var invs []models.Inventory
	if err := json.Unmarshal(getResp.Body(), &invs); err != nil || len(invs) == 0 {
		h.logger.Warnf("bad inventory response: %v", err)
		return nil
	}
	available := invs[0].QuantityOnHand
	if available < oc.Quantity {
		h.logger.Warnf("insufficient stock for %s: have=%d need=%d",
			oc.OrderID, available, oc.Quantity)
		return nil
	}

	// 3) PUT to deduct
	putReq := &commonsHttp.Request{
		Url: "/inventory",
		Body: map[string]interface{}{
			"tenant_id": oc.TenantID,
			"hub_id":    oc.HubID,
			"sku_id":    oc.SKUID,
			"quantity":  available - oc.Quantity,
		},
		Timeout: 5 * time.Second,
	}
	if _, err := h.client.Put(putReq, nil); err != nil {
		h.logger.Errorf("IMS PUT error: %v", err)
		return err
	}

	// 4) Update Mongo
	if _, err := h.coll.UpdateOne(ctx,
		bson.M{"_id": oc.OrderID},
		bson.M{"$set": bson.M{"status": "new_order", "updated_at": time.Now().UTC()}},
	); err != nil {
		h.logger.Errorf("mongo update failed: %v", err)
		return err
	}
	h.logger.Infof("Order %s → new_order", oc.OrderID)

	// 5) Publish order.updated
	evt, _ := json.Marshal(oc)
	if err := h.publisher.Publish(ctx, &pubsub.Message{
		Topic: updatedTopic,
		Key:   oc.OrderID,
		Value: evt,
	}); err != nil {
		h.logger.Errorf("publish updated failed: %v", err)
		return err
	}
	h.logger.Infof("Published order.updated for %s", oc.OrderID)
	return nil
}
