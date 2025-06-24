package webhooklogs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	coll   *mongo.Collection
	logger *log.Logger
}

func NewHandler(coll *mongo.Collection) *Handler {
	return &Handler{coll: coll, logger: log.DefaultLogger()}
}

func (h *Handler) Process(ctx context.Context, msg *pubsub.Message) error {
	var failed struct {
		WebhookID   string                 `json:"webhookId"`
		Event       string                 `json:"event"`
		CallbackURL string                 `json:"callbackUrl"`
		Error       string                 `json:"error"`
		Payload     map[string]interface{} `json:"payload"`
	}
	if err := json.Unmarshal(msg.Value, &failed); err != nil {
		h.logger.Errorf("webhook_logs: invalid JSON: %v", err)
		return err
	}

	record := FailedWebhook{
		WebhookID:   failed.WebhookID,
		Event:       failed.Event,
		CallbackURL: failed.CallbackURL,
		Error:       failed.Error,
		Payload:     failed.Payload,
		CreatedAt:   time.Now().UTC(),
	}
	if _, err := h.coll.InsertOne(ctx, record); err != nil {
		h.logger.Errorf("webhook_logs: mongo insert failed: %v", err)
		return err
	}
	h.logger.Infof("webhook_logs: recorded failure %s", record.WebhookID)
	return nil
}
