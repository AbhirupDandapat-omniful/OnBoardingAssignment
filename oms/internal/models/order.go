package models

import (
	"strconv"
	"time"
)

type Order struct {
	ID        string    `bson:"_id,omitempty"`
	TenantID  string    `bson:"tenant_id"`
	SellerID  string    `bson:"seller_id"`
	HubID     string    `bson:"hub_id"`
	SKUID     string    `bson:"sku_id"`
	Quantity  int64     `bson:"quantity"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
}

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

	idx := make(map[string]int, len(header))
	for i, col := range header {
		idx[col] = i
	}

	qty, err := strconv.ParseInt(rec[idx["quantity"]], 10, 64)
	if err != nil || qty <= 0 {
		return Order{}, false
	}

	t := rec[idx["tenant_id"]]
	s := rec[idx["seller_id"]]
	h := rec[idx["hub_id"]]
	k := rec[idx["sku_id"]]
	if t == "" || s == "" || h == "" || k == "" {
		return Order{}, false
	}

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
