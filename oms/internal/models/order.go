package models

import (
	"strconv"
	"time"
)

// Order represents a row from your CSV persisted to MongoDB.
type Order struct {
	ID        string    `bson:"_id,omitempty"` // Mongo document ID
	TenantID  string    `bson:"tenant_id"`
	SellerID  string    `bson:"seller_id"`
	HubID     string    `bson:"hub_id"`
	SKUID     string    `bson:"sku_id"`
	Quantity  int64     `bson:"quantity"`
	Status    string    `bson:"status"` // e.g. "on_hold", "new_order"
	CreatedAt time.Time `bson:"created_at"`
}

// OrderCreated is the Kafka event payload you’ll emit.
type OrderCreated struct {
	OrderID   string    `json:"order_id"`
	TenantID  string    `json:"tenant_id"`
	SellerID  string    `json:"seller_id"`
	HubID     string    `json:"hub_id"`
	SKUID     string    `json:"sku_id"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

func OrderFromRow(rec []string, header []string) (Order, bool) {
	// build a map: column name → its index
	idx := make(map[string]int, len(header))
	for i, col := range header {
		idx[col] = i
	}

	// parse quantity (must be > 0)
	qty, err := strconv.ParseInt(rec[idx["quantity"]], 10, 64)
	if err != nil || qty <= 0 {
		return Order{}, false
	}

	// extract the four IDs
	t := rec[idx["tenant_id"]]
	s := rec[idx["seller_id"]]
	h := rec[idx["hub_id"]]
	k := rec[idx["sku_id"]]
	if t == "" || s == "" || h == "" || k == "" {
		return Order{}, false
	}

	// build the Order
	return Order{
		TenantID:  t,
		SellerID:  s,
		HubID:     h,
		SKUID:     k,
		Quantity:  qty,
		Status:    "on_hold",
		CreatedAt: time.Now().UTC(),
	}, true
}
