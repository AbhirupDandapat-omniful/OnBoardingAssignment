package webhooklogs

import "time"

type FailedWebhook struct {
	ID          string                 `bson:"_id,omitempty"`
	WebhookID   string                 `bson:"webhook_id"`
	Event       string                 `bson:"event"`
	CallbackURL string                 `bson:"callback_url"`
	Error       string                 `bson:"error"`
	Payload     map[string]interface{} `bson:"payload"`
	CreatedAt   time.Time              `bson:"created_at"`
}
