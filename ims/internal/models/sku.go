package models

import "time"

type SKU struct {
    ID          string    `db:"id"            json:"id"`
    TenantID    string    `db:"tenant_id"     json:"tenant_id"`
    SellerID    string    `db:"seller_id"     json:"seller_id"`
    Code        string    `db:"code"          json:"code"`
    Name        string    `db:"name"          json:"name"`
    Description string    `db:"description"   json:"description,omitempty"`
    CategoryID  string    `db:"category_id"   json:"category_id,omitempty"`
    Weight      float64   `db:"weight"        json:"weight,omitempty"`
    WeightUnit  string    `db:"weight_unit"   json:"weight_unit,omitempty"`
    Length      float64   `db:"length"        json:"length,omitempty"`
    Width       float64   `db:"width"         json:"width,omitempty"`
    Height      float64   `db:"height"        json:"height,omitempty"`
    CreatedAt   time.Time `db:"created_at"    json:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"    json:"updated_at"`
}
