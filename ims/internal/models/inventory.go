package models

import "time"

type Inventory struct {
	HubID            string    `json:"hub_id"            gorm:"column:hub_id"`
	SKUID            string    `json:"sku_id"            gorm:"column:sku_id"`
	QuantityOnHand   int64     `json:"quantity_on_hand"  gorm:"column:quantity_on_hand"`
	QuantityReserved int64     `json:"quantity_reserved" gorm:"column:quantity_reserved"`
	MinThreshold     int64     `json:"min_threshold"     gorm:"column:min_threshold"`
	MaxThreshold     int64     `json:"max_threshold"     gorm:"column:max_threshold"`
	UpdatedAt        time.Time `json:"updated_at"        gorm:"column:updated_at"`
}
