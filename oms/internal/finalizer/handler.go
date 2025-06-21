// internal/finalizer/handler.go
package finalizer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/omniful/go_commons/config"
	commonsHttp "github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/abhirup.dandapat/oms/internal/models"
)

type Handler struct {
	logger *log.Logger
	client *commonsHttp.Client
	coll   *mongo.Collection
}

func NewHandler() *Handler {
	ctx, _ := config.TODOContext()
	logger := log.DefaultLogger()

	// 1) Build a non-nil HTTP transport
	transport := &http.Transport{}

	// 2) Pass it into NewHTTPClient so .Transport is never nil
	baseURL := config.GetString(ctx, "ims.baseUrl")
	httpClient, err := commonsHttp.NewHTTPClient("oms-finalizer", baseURL, transport)
	if err != nil {
		logger.Panicf("creating HTTP client: %v", err)
	}

	// 3) Mongo setup (unchanged)
	mongoURI := config.GetString(ctx, "mongo.uri")
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Panicf("connecting to Mongo: %v", err)
	}
	coll := mongoClient.Database("omsdb").Collection("orders")

	return &Handler{
		logger: logger,
		client: httpClient,
		coll:   coll,
	}
}

// Process satisfies pubsub.IPubSubMessageHandler
func (h *Handler) Process(ctx context.Context, msg *pubsub.Message) error {
	// 1) decode the OrderCreated event
	var evt models.OrderCreated
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		h.logger.Errorf("unmarshal OrderCreated: %v", err)
		return err
	}
	h.logger.Infof("Finalizing order %s", evt.OrderID)

	// 2) fetch inventory from IMS
	path := fmt.Sprintf("/inventory?hub_id=%s&sku_ids=%s", evt.HubID, evt.SKUID)
	resp, err := h.client.Get(&commonsHttp.Request{Url: path}, nil)
	if err != nil {
		h.logger.Errorf("inventory lookup failed: %v", err)
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		h.logger.Warnf("inventory lookup returned status %d", resp.StatusCode())
		return nil
	}

	// 3) parse the JSON array
	var invs []models.Inventory
	if err := json.Unmarshal(resp.Body(), &invs); err != nil {
		h.logger.Errorf("unmarshal inventory: %v", err)
		return err
	}
	if len(invs) == 0 {
		h.logger.Warn("inventory response empty")
		return nil
	}

	inv := invs[0]
	available := inv.QuantityOnHand - inv.QuantityReserved

	// 4) if enough stock → upsert new quantity & update order
	if available >= evt.Quantity {
		// a) reserve in IMS
		up := struct {
			TenantID string `json:"tenant_id"`
			HubID    string `json:"hub_id"`
			SKUID    string `json:"sku_id"`
			Quantity int64  `json:"quantity"`
		}{
			evt.TenantID, evt.HubID, evt.SKUID, inv.QuantityOnHand - evt.Quantity,
		}
		if _, err := h.client.Put(&commonsHttp.Request{
			Url:     "/inventory",
			Body:    up,
			Timeout: 5 * time.Second,
		}, nil); err != nil {
			h.logger.Errorf("inventory upsert failed: %v", err)
			return err
		}

		// b) update Mongo
		if _, err := h.coll.UpdateOne(ctx,
			bson.M{"_id": evt.OrderID},
			bson.M{"$set": bson.M{"status": "new_order"}},
		); err != nil {
			h.logger.Errorf("Mongo update status: %v", err)
			return err
		}
		h.logger.Infof("Order %s → new_order", evt.OrderID)
	} else {
		h.logger.Infof("Order %s remains on_hold (%d available)", evt.OrderID, available)
	}

	return nil
}
