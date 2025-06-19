package models

import "time"

// Category classifies SKUs.
type Category struct {
    ID          string    `db:"id"          json:"id"`
    TenantID    string    `db:"tenant_id"   json:"tenant_id"`
    Name        string    `db:"name"        json:"name"`
    Description string    `db:"description" json:"description,omitempty"`
    CreatedAt   time.Time `db:"created_at"  json:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
}
