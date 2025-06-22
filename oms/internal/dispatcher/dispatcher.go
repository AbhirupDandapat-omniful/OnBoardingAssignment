package dispatcher

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	commonsHttp "github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Webhook struct {
	ID          string            `bson:"_id"          json:"id"`
	TenantID    string            `bson:"tenant_id"    json:"tenant_id"`
	CallbackURL string            `bson:"callback_url" json:"callback_url"`
	Events      []string          `bson:"events"       json:"events"`
	Headers     map[string]string `bson:"headers"      json:"headers"`
	IsActive    bool              `bson:"is_active"    json:"is_active"`
}

type Dispatcher struct {
	coll        *mongo.Collection
	httpClient  *commonsHttp.Client
	producer    pubsub.Publisher
	failedTopic string
	logger      *log.Logger
}

func NewDispatcher(
	coll *mongo.Collection,
	httpClient *commonsHttp.Client,
	producer pubsub.Publisher,
	failedTopic string,
) *Dispatcher {
	return &Dispatcher{
		coll:        coll,
		httpClient:  httpClient,
		producer:    producer,
		failedTopic: failedTopic,
		logger:      log.DefaultLogger(),
	}
}

// Process implements pubsub.IPubSubMessageHandler
func (d *Dispatcher) Process(ctx context.Context, msg *pubsub.Message) error {
	eventType := msg.Topic

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		d.logger.Errorf("invalid JSON payload: %v", err)
		return err
	}

	// ← key must be "tenant_id"
	tenantID, ok := payload["tenant_id"].(string)
	if !ok {
		d.logger.Errorf("missing tenant_id in payload")
		return nil
	}

	filter := bson.M{
		"tenant_id": tenantID,
		"is_active": true,
		"events":    eventType,
	}
	cursor, err := d.coll.Find(ctx, filter)
	if err != nil {
		d.logger.Errorf("mongo Find webhooks: %v", err)
		return err
	}
	defer cursor.Close(ctx)

	var whs []Webhook
	if err := cursor.All(ctx, &whs); err != nil {
		d.logger.Errorf("cursor.All webhooks: %v", err)
		return err
	}
	if len(whs) == 0 {
		return nil
	}

	body := map[string]interface{}{
		"event": eventType,
		"data":  payload,
	}

	for _, w := range whs {
		hdrs := make(http.Header)
		for k, v := range w.Headers {
			hdrs.Add(k, v)
		}
		req := &commonsHttp.Request{
			Url:     w.CallbackURL,
			Body:    body,
			Timeout: 5 * time.Second,
			Headers: hdrs,
		}
		if _, err := d.httpClient.Post(req, nil); err != nil {
			d.logger.Warnf("webhook POST failed (%s → %s): %v", w.ID, w.CallbackURL, err)
			dead := map[string]interface{}{
				"webhookId":   w.ID,
				"event":       eventType,
				"callbackUrl": w.CallbackURL,
				"error":       err.Error(),
				"payload":     payload,
			}
			if buf, jerr := json.Marshal(dead); jerr == nil {
				d.producer.Publish(ctx, &pubsub.Message{
					Topic: d.failedTopic,
					Key:   w.ID,
					Value: buf,
				})
			}
			continue
		}
		d.logger.Infof("webhook delivered: %s → %s", w.ID, w.CallbackURL)
	}

	return nil
}
