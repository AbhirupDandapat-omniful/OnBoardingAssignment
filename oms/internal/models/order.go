package models

import "time"

type Order struct {
    TenantID  string    `bson:"tenant_id"`  
    SellerID  string    `bson:"seller_id"`
    HubID     string    `bson:"hub_id"`
    SKUID     string    `bson:"sku_id"`
    Quantity  int       `bson:"quantity"`
    Status    string    `bson:"status"`    
    CreatedAt time.Time `bson:"created_at"` 
}
