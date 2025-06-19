package models

import "time"

// InventoryTransaction logs every stock movement.
type InventoryTransaction struct {
    ID              string    `db:"id"              json:"id"`
    TenantID        string    `db:"tenant_id"       json:"tenant_id"`
    HubID           string    `db:"hub_id"          json:"hub_id"`
    SKUID           string    `db:"sku_id"          json:"sku_id"`
    Delta           int64     `db:"delta"           json:"delta"`
    TransactionType string    `db:"transaction_type"json:"transaction_type"`
    ReferenceID     string    `db:"reference_id"    json:"reference_id,omitempty"`
    CreatedAt       time.Time `db:"created_at"      json:"created_at"`
}
