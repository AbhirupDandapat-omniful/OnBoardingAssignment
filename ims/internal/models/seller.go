package models

import "time"

// Seller represents a seller under a tenant.
type Seller struct {
    ID        string                 `db:"id"         json:"id"`
    TenantID  string                 `db:"tenant_id"  json:"tenant_id"`
    Name      string                 `db:"name"       json:"name"`
    Metadata  map[string]interface{} `db:"metadata"   json:"metadata,omitempty"`
    CreatedAt time.Time              `db:"created_at" json:"created_at"`
    UpdatedAt time.Time              `db:"updated_at" json:"updated_at"`
}
