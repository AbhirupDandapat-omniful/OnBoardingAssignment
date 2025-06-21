package models

import "time"

// Inventory mirrors the IMS /inventory response.
type Inventory struct {
	HubID            string    `json:"hub_id" bson:"hub_id"`
	SKUID            string    `json:"sku_id" bson:"sku_id"`
	QuantityOnHand   int64     `json:"quantity_on_hand" bson:"quantity_on_hand"`
	QuantityReserved int64     `json:"quantity_reserved" bson:"quantity_reserved"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}
