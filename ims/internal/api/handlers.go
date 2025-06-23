package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"

	"github.com/abhirup.dandapat/ims/internal/models"
	"github.com/abhirup.dandapat/ims/internal/store"
)

func createTenant(c *gin.Context) {
	var t models.Tenant
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	t.ID = uuid.New().String()
	now := time.Now().UTC()
	t.CreatedAt, t.UpdatedAt = now, now

	metaBytes, err := json.Marshal(t.Metadata)
	if err != nil {
		log.DefaultLogger().Errorf("createTenant: failed to marshal metadata: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`INSERT INTO tenants(id,name,metadata,created_at,updated_at)
           VALUES(?,?,?,?,?)`,
		t.ID, t.Name, metaBytes, t.CreatedAt, t.UpdatedAt,
	).Error; err != nil {
		log.DefaultLogger().Errorf("createTenant DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_tenant_failed")})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func getTenant(c *gin.Context) {
	id := c.Param("id")
	var t models.Tenant

	db := store.DB.GetSlaveDB(c.Request.Context())
	if err := db.Raw(
		`SELECT id,name,metadata,created_at,updated_at
         FROM tenants WHERE id = ?`, id,
	).Scan(&t).Error; err != nil {
		log.DefaultLogger().Errorf("getTenant DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.tenant_not_found")})
		return
	}

	c.JSON(http.StatusOK, t)
}

func createSeller(c *gin.Context) {
	var s models.Seller
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	s.ID = uuid.New().String()
	now := time.Now().UTC()
	s.CreatedAt, s.UpdatedAt = now, now

	metaBytes, err := json.Marshal(s.Metadata)
	if err != nil {
		log.DefaultLogger().Errorf("createSeller: failed to marshal metadata: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`INSERT INTO sellers(id,tenant_id,name,metadata,created_at,updated_at)
           VALUES(?,?,?,?,?,?)`,
		s.ID, s.TenantID, s.Name, metaBytes, s.CreatedAt, s.UpdatedAt,
	).Error; err != nil {
		log.DefaultLogger().Errorf("createSeller DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_seller_failed")})
		return
	}

	c.JSON(http.StatusCreated, s)
}
func getSeller(c *gin.Context) {
	id := c.Param("id")
	var s models.Seller

	db := store.DB.GetSlaveDB(c.Request.Context())
	if err := db.Raw(
		`SELECT id,tenant_id,name,metadata,created_at,updated_at
         FROM sellers WHERE id = ?`, id,
	).Scan(&s).Error; err != nil {
		log.DefaultLogger().Errorf("getSeller DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.seller_not_found")})
		return
	}

	c.JSON(http.StatusOK, s)
}

func createCategory(c *gin.Context) {
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	cat.ID = uuid.New().String()
	now := time.Now().UTC()
	cat.CreatedAt, cat.UpdatedAt = now, now

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`INSERT INTO categories(id,tenant_id,name,description,created_at,updated_at)
         VALUES(?,?,?,?,?,?)`,
		cat.ID, cat.TenantID, cat.Name, cat.Description, cat.CreatedAt, cat.UpdatedAt,
	).Error; err != nil {
		log.DefaultLogger().Errorf("createCategory DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_category_failed")})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

func getCategory(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category

	db := store.DB.GetSlaveDB(c.Request.Context())
	if err := db.Raw(
		`SELECT id,tenant_id,name,description,created_at,updated_at
         FROM categories WHERE id = ?`, id,
	).Scan(&cat).Error; err != nil {
		log.DefaultLogger().Errorf("getCategory DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.category_not_found")})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func createHub(c *gin.Context) {
	var h models.Hub
	if err := c.ShouldBindJSON(&h); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	h.ID = uuid.New().String()
	now := time.Now().UTC()
	h.CreatedAt, h.UpdatedAt = now, now

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`INSERT INTO hubs(id,tenant_id,seller_id,name,location,address,contact_email,contact_phone,timezone,created_at,updated_at)
         VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
		h.ID, h.TenantID, h.SellerID, h.Name, h.Location,
		h.Address, h.ContactEmail, h.ContactPhone, h.Timezone,
		h.CreatedAt, h.UpdatedAt,
	).Error; err != nil {
		log.DefaultLogger().Errorf("createHub DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_hub_failed")})
		return
	}

	c.JSON(http.StatusCreated, h)
}

func getHub(c *gin.Context) {
	id := c.Param("id")

	if cached, err := store.RedisClient.Get(c.Request.Context(), "hub:"+id); err == nil {
		c.JSON(http.StatusOK, cached)
		return
	}

	var h models.Hub
	db := store.DB.GetSlaveDB(c.Request.Context())
	if err := db.Raw(
		`SELECT id,tenant_id,seller_id,name,location,address,contact_email,contact_phone,timezone,created_at,updated_at
         FROM hubs WHERE id = ?`, id,
	).Scan(&h).Error; err != nil {
		log.DefaultLogger().Errorf("getHub DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.hub_not_found")})
		return
	}

	if b, err := json.Marshal(h); err == nil {
		_, _ = store.RedisClient.Set(c.Request.Context(), "hub:"+id, string(b), 5*time.Minute)
	}

	c.JSON(http.StatusOK, h)
}

func updateHub(c *gin.Context) {
	id := c.Param("id")
	var h models.Hub
	if err := c.ShouldBindJSON(&h); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	h.UpdatedAt = time.Now().UTC()

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`UPDATE hubs SET name=?,location=?,address=?,contact_email=?,contact_phone=?,timezone=?,updated_at=? WHERE id=?`,
		h.Name, h.Location, h.Address, h.ContactEmail, h.ContactPhone, h.Timezone, h.UpdatedAt, id,
	).Error; err != nil {
		log.DefaultLogger().Errorf("updateHub DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_hub_failed")})
		return
	}

	_, _ = store.RedisClient.Del(c.Request.Context(), "hub:"+id)
	c.Status(http.StatusOK)
}

func deleteHub(c *gin.Context) {
	id := c.Param("id")

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(`DELETE FROM hubs WHERE id = ?`, id).Error; err != nil {
		log.DefaultLogger().Errorf("deleteHub DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_hub_failed")})
		return
	}

	_, _ = store.RedisClient.Del(c.Request.Context(), "hub:"+id)
	c.Status(http.StatusNoContent)
}

func listHubs(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	sellerID := c.Query("seller_id")

	where := []string{"1=1"}
	args := []interface{}{}
	if tenantID != "" {
		where = append(where, "tenant_id = ?")
		args = append(args, tenantID)
	}
	if sellerID != "" {
		where = append(where, "seller_id = ?")
		args = append(args, sellerID)
	}

	sqlStr := fmt.Sprintf(
		`SELECT id,tenant_id,seller_id,name,location,address,contact_email,contact_phone,timezone,created_at,updated_at
           FROM hubs WHERE %s`, strings.Join(where, " AND "),
	)

	db := store.DB.GetSlaveDB(c.Request.Context())
	rows, err := db.Raw(sqlStr, args...).Rows()
	if err != nil {
		log.DefaultLogger().Errorf("listHubs DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.list_hubs_failed")})
		return
	}
	defer rows.Close()

	var hubs []models.Hub
	for rows.Next() {
		var h models.Hub
		if err := db.ScanRows(rows, &h); err != nil {
			log.DefaultLogger().Errorf("listHubs scan error: %v", err)
			continue
		}
		hubs = append(hubs, h)
	}
	c.JSON(http.StatusOK, hubs)
}

func createSKU(c *gin.Context) {
	var s models.SKU
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	s.ID = uuid.New().String()
	now := time.Now().UTC()
	s.CreatedAt, s.UpdatedAt = now, now

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`INSERT INTO skus(id,tenant_id,seller_id,code,name,description,category_id,weight,weight_unit,length,width,height,created_at,updated_at)
         VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		s.ID, s.TenantID, s.SellerID, s.Code, s.Name, s.Description,
		s.CategoryID, s.Weight, s.WeightUnit, s.Length, s.Width, s.Height,
		s.CreatedAt, s.UpdatedAt,
	).Error; err != nil {
		log.DefaultLogger().Errorf("createSKU DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_sku_failed")})
		return
	}

	c.JSON(http.StatusCreated, s)
}

func getSKU(c *gin.Context) {
	id := c.Param("id")

	if cached, err := store.RedisClient.Get(c.Request.Context(), "sku:"+id); err == nil {
		c.JSON(http.StatusOK, cached)
		return
	}

	var s models.SKU
	db := store.DB.GetSlaveDB(c.Request.Context())
	if err := db.Raw(
		`SELECT id,tenant_id,seller_id,code,name,description,category_id,weight,weight_unit,length,width,height,created_at,updated_at
         FROM skus WHERE id = ?`, id,
	).Scan(&s).Error; err != nil {
		log.DefaultLogger().Errorf("getSKU DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.sku_not_found")})
		return
	}

	if b, err := json.Marshal(s); err == nil {
		_, _ = store.RedisClient.Set(c.Request.Context(), "sku:"+id, string(b), 5*time.Minute)
	}

	c.JSON(http.StatusOK, s)
}

func updateSKU(c *gin.Context) {
	id := c.Param("id")
	var s models.SKU
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	s.UpdatedAt = time.Now().UTC()

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`UPDATE skus SET code=?,name=?,description=?,category_id=?,weight=?,weight_unit=?,length=?,width=?,height=?,updated_at=? WHERE id=?`,
		s.Code, s.Name, s.Description, s.CategoryID,
		s.Weight, s.WeightUnit, s.Length, s.Width, s.Height,
		s.UpdatedAt, id,
	).Error; err != nil {
		log.DefaultLogger().Errorf("updateSKU DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_sku_failed")})
		return
	}

	_, _ = store.RedisClient.Del(c.Request.Context(), "sku:"+id)
	c.Status(http.StatusOK)
}

func deleteSKU(c *gin.Context) {
	id := c.Param("id")

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(`DELETE FROM skus WHERE id = ?`, id).Error; err != nil {
		log.DefaultLogger().Errorf("deleteSKU DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_sku_failed")})
		return
	}

	_, _ = store.RedisClient.Del(c.Request.Context(), "sku:"+id)
	c.Status(http.StatusNoContent)
}

func listSKUs(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	codesParam := c.Query("sku_codes")
	codes := []string{}
	if codesParam != "" {
		codes = strings.Split(codesParam, ",")
	}

	where := []string{"1=1"}
	args := []interface{}{}
	if tenantID != "" {
		where = append(where, "tenant_id = ?")
		args = append(args, tenantID)
	}
	if len(codes) > 0 {
		ph := make([]string, len(codes))
		for i := range codes {
			ph[i] = "?"
		}
		where = append(where, fmt.Sprintf("code IN (%s)", strings.Join(ph, ",")))
		for _, code := range codes {
			args = append(args, code)
		}
	}

	sqlStr := fmt.Sprintf(
		`SELECT id,tenant_id,seller_id,code,name,description,category_id,weight,weight_unit,length,width,height,created_at,updated_at
           FROM skus WHERE %s`, strings.Join(where, " AND "),
	)

	db := store.DB.GetSlaveDB(c.Request.Context())
	rows, err := db.Raw(sqlStr, args...).Rows()
	if err != nil {
		log.DefaultLogger().Errorf("listSKUs DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.list_skus_failed")})
		return
	}
	defer rows.Close()

	var skus []models.SKU
	for rows.Next() {
		var s models.SKU
		if err := db.ScanRows(rows, &s); err != nil {
			log.DefaultLogger().Errorf("listSKUs scan error: %v", err)
			continue
		}
		skus = append(skus, s)
	}
	c.JSON(http.StatusOK, skus)
}

type InventoryUpsertRequest struct {
	TenantID string `json:"tenant_id" binding:"required"`
	HubID    string `json:"hub_id"    binding:"required"`
	SKUID    string `json:"sku_id"    binding:"required"`
	Quantity int64  `json:"quantity"  binding:"required"`
}

func upsertInventory(c *gin.Context) {
	var req InventoryUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	now := time.Now().UTC()

	tx := store.DB.GetMasterDB(c.Request.Context()).Begin()
	if tx.Error != nil {
		log.DefaultLogger().Errorf("upsertInventory begin tx error: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.inventory_upsert_failed")})
		return
	}
	defer tx.Rollback()

	if err := tx.Exec(
		`INSERT INTO inventory(hub_id,sku_id,quantity_on_hand,quantity_reserved,min_threshold,max_threshold,updated_at)
         VALUES(?,?,?,?,?,?,?)
         ON CONFLICT (hub_id,sku_id) DO UPDATE SET quantity_on_hand = EXCLUDED.quantity_on_hand, updated_at = EXCLUDED.updated_at`,
		req.HubID, req.SKUID, req.Quantity, 0, 0, 0, now,
	).Error; err != nil {
		log.DefaultLogger().Errorf("upsertInventory exec error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.inventory_upsert_failed")})
		return
	}

	transID := uuid.New().String()
	if err := tx.Exec(
		`INSERT INTO inventory_transactions(id,tenant_id,hub_id,sku_id,delta,transaction_type,reference_id,created_at)
         VALUES(?,?,?,?,?,?,?,?)`,
		transID, req.TenantID, req.HubID, req.SKUID, req.Quantity, "upsert", "", now,
	).Error; err != nil {
		log.DefaultLogger().Errorf("upsertInventory trans error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.inventory_upsert_failed")})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.DefaultLogger().Errorf("upsertInventory commit error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.inventory_upsert_failed")})
		return
	}

	var inv models.Inventory
	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Raw(
		`SELECT hub_id,sku_id,quantity_on_hand,quantity_reserved,min_threshold,max_threshold,updated_at
         FROM inventory WHERE hub_id = ? AND sku_id = ?`,
		req.HubID, req.SKUID,
	).Scan(&inv).Error; err != nil {
		log.DefaultLogger().Errorf("upsertInventory fetch error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.inventory_upsert_failed")})
		return
	}

	c.JSON(http.StatusOK, inv)
}

func listInventory(c *gin.Context) {
	hubID := c.Query("hub_id")
	skuIDs := strings.Split(c.Query("sku_ids"), ",")

	where := []string{"hub_id = ?"}
	args := []interface{}{hubID}
	if len(skuIDs) > 0 && skuIDs[0] != "" {
		ph := strings.Repeat("?,", len(skuIDs))
		ph = ph[:len(ph)-1]
		where = append(where, fmt.Sprintf("sku_id IN (%s)", ph))
		for _, id := range skuIDs {
			args = append(args, id)
		}
	}

	sqlStr := fmt.Sprintf(
		`SELECT hub_id,sku_id,quantity_on_hand,quantity_reserved,min_threshold,max_threshold,updated_at
           FROM inventory WHERE %s`, strings.Join(where, " AND "),
	)

	db := store.DB.GetSlaveDB(c.Request.Context())
	rows, err := db.Raw(sqlStr, args...).Rows()
	if err != nil {
		log.DefaultLogger().Errorf("listInventory DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.list_inventory_failed")})
		return
	}
	defer rows.Close()

	var invs []models.Inventory
	for rows.Next() {
		var inv models.Inventory
		if err := db.ScanRows(rows, &inv); err != nil {
			log.DefaultLogger().Errorf("listInventory scan error: %v", err)
			continue
		}
		invs = append(invs, inv)
	}

	c.JSON(http.StatusOK, invs)
}

func createWebhook(c *gin.Context) {
	var w models.WebhookRegistration
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	w.ID = uuid.New().String()
	now := time.Now().UTC()
	w.CreatedAt, w.UpdatedAt = now, now

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`INSERT INTO webhooks(id,tenant_id,callback_url,events,headers,is_active,created_at,updated_at)
         VALUES(?,?,?,?,?,?,?,?)`,
		w.ID, w.TenantID, w.CallbackURL, w.Events, w.Headers, w.IsActive, w.CreatedAt, w.UpdatedAt,
	).Error; err != nil {
		log.DefaultLogger().Errorf("createWebhook DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_webhook_failed")})
		return
	}

	c.JSON(http.StatusCreated, w)
}

func getWebhook(c *gin.Context) {
	id := c.Param("id")
	var w models.WebhookRegistration

	db := store.DB.GetSlaveDB(c.Request.Context())
	if err := db.Raw(
		`SELECT id,tenant_id,callback_url,events,headers,is_active,created_at,updated_at
         FROM webhooks WHERE id = ?`, id,
	).Scan(&w).Error; err != nil {
		log.DefaultLogger().Errorf("getWebhook DB error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.webhook_not_found")})
		return
	}

	c.JSON(http.StatusOK, w)
}

func updateWebhook(c *gin.Context) {
	id := c.Param("id")
	var w models.WebhookRegistration
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	w.UpdatedAt = time.Now().UTC()

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`UPDATE webhooks
         SET callback_url=?,events=?,headers=?,is_active=?,updated_at=?
         WHERE id=?`,
		w.CallbackURL, w.Events, w.Headers, w.IsActive, w.UpdatedAt, id,
	).Error; err != nil {
		log.DefaultLogger().Errorf("updateWebhook DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_webhook_failed")})
		return
	}

	c.Status(http.StatusOK)
}

func deleteWebhook(c *gin.Context) {
	id := c.Param("id")

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(`DELETE FROM webhooks WHERE id = ?`, id).Error; err != nil {
		log.DefaultLogger().Errorf("deleteWebhook DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_webhook_failed")})
		return
	}

	c.Status(http.StatusNoContent)
}
