package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/abhirup.dandapat/ims/internal/models"
	"github.com/abhirup.dandapat/ims/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

var invTxLogger = log.DefaultLogger()

type InventoryTransactionRequest struct {
	TenantID        string `json:"tenant_id"        form:"tenant_id"`
	HubID           string `json:"hub_id"           form:"hub_id"`
	SKUID           string `json:"sku_id"           form:"sku_id"`
	Delta           int64  `json:"delta"            form:"delta"`
	TransactionType string `json:"transaction_type" form:"transaction_type"`
	ReferenceID     string `json:"reference_id"     form:"reference_id"`
}

func createInventoryTransaction(c *gin.Context) {
	var req InventoryTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	tx := models.InventoryTransaction{
		ID:              uuid.New().String(),
		TenantID:        req.TenantID,
		HubID:           req.HubID,
		SKUID:           req.SKUID,
		Delta:           req.Delta,
		TransactionType: req.TransactionType,
		ReferenceID:     req.ReferenceID,
		CreatedAt:       time.Now().UTC(),
	}

	db := store.DB.GetMasterDB(c.Request.Context())
	if err := db.Exec(
		`INSERT INTO inventory_transactions
		 (id,tenant_id,hub_id,sku_id,delta,transaction_type,reference_id,created_at)
		 VALUES(?,?,?,?,?,?,?,?)`,
		tx.ID, tx.TenantID, tx.HubID, tx.SKUID,
		tx.Delta, tx.TransactionType, tx.ReferenceID, tx.CreatedAt,
	).Error; err != nil {
		invTxLogger.Errorf("createInventoryTransaction DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.inventory_transaction_failed")})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func listInventoryTransactions(c *gin.Context) {
	var req InventoryTransactionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	where := []string{"1=1"}
	args := []interface{}{}
	if req.TenantID != "" {
		where = append(where, "tenant_id = ?")
		args = append(args, req.TenantID)
	}
	if req.HubID != "" {
		where = append(where, "hub_id = ?")
		args = append(args, req.HubID)
	}
	if req.SKUID != "" {
		where = append(where, "sku_id = ?")
		args = append(args, req.SKUID)
	}

	sql := `SELECT id,tenant_id,hub_id,sku_id,delta,transaction_type,reference_id,created_at
	        FROM inventory_transactions
	        WHERE ` + strings.Join(where, " AND ") + `
	        ORDER BY created_at DESC`

	db := store.DB.GetSlaveDB(c.Request.Context())
	rows, err := db.Raw(sql, args...).Rows()
	if err != nil {
		invTxLogger.Errorf("listInventoryTransactions DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.inventory_transaction_list_failed")})
		return
	}
	defer rows.Close()

	var txs []models.InventoryTransaction
	for rows.Next() {
		var t models.InventoryTransaction
		if err := db.ScanRows(rows, &t); err != nil {
			invTxLogger.Warnf("scan inventory_transaction row: %v", err)
			continue
		}
		txs = append(txs, t)
	}

	c.JSON(http.StatusOK, gin.H{"transactions": txs})
}
