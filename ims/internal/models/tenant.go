package models

import "time"

type Tenant struct {
	ID        string                 `db:"id"         json:"id"`
	Name      string                 `db:"name"       json:"name"`
	Metadata  map[string]interface{} `db:"metadata"   json:"metadata,omitempty"`
	CreatedAt time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt time.Time              `db:"updated_at" json:"updated_at"`
}
