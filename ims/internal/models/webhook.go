package models

import "time"

type WebhookRegistration struct {
	ID          string            `db:"id"           json:"id"`
	TenantID    string            `db:"tenant_id"    json:"tenant_id"`
	CallbackURL string            `db:"callback_url" json:"callback_url"`
	Events      []string          `db:"events"       json:"events"`
	Headers     map[string]string `db:"headers"      json:"headers,omitempty"`
	IsActive    bool              `db:"is_active"    json:"is_active"`
	CreatedAt   time.Time         `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at"   json:"updated_at"`
}
