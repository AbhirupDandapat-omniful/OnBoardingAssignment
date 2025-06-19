package models

import "time"

type Inventory struct {
    HubID            string    `db:"hub_id"              json:"hub_id"`
    SKUID            string    `db:"sku_id"              json:"sku_id"`
    QuantityOnHand   int64     `db:"quantity"            json:"quantity"`
    QuantityReserved int64     `db:"quantity_reserved"   json:"quantity_reserved"`
    MinThreshold     int64     `db:"min_threshold"       json:"min_threshold"`
    MaxThreshold     int64     `db:"max_threshold"       json:"max_threshold"`
    UpdatedAt        time.Time `db:"updated_at"          json:"updated_at"`
}
