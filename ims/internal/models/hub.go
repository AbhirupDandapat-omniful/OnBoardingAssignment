package models

import "time"

type Hub struct {
    ID           string    `db:"id"            json:"id"`
    TenantID     string    `db:"tenant_id"     json:"tenant_id"`
    SellerID     string    `db:"seller_id"     json:"seller_id"`
    Name         string    `db:"name"          json:"name"`
    Location     string    `db:"location"      json:"location"`
    Address      string    `db:"address"       json:"address,omitempty"`
    ContactEmail string    `db:"contact_email" json:"contact_email,omitempty"`
    ContactPhone string    `db:"contact_phone" json:"contact_phone,omitempty"`
    Timezone     string    `db:"timezone"      json:"timezone,omitempty"`
    CreatedAt    time.Time `db:"created_at"    json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}
