package models

import "time"

type Webhook struct {
	ID          string            `bson:"_id"           json:"id"`
	TenantID    string            `bson:"tenant_id"     json:"tenant_id"`
	CallbackURL string            `bson:"callback_url"  json:"callback_url"`
	Events      []string          `bson:"events"        json:"events"`
	Headers     map[string]string `bson:"headers"       json:"headers"`
	IsActive    bool              `bson:"is_active"     json:"is_active"`
	CreatedAt   time.Time         `bson:"created_at"    json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at"    json:"updated_at"`
}
