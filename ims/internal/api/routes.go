package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {

	r.POST("/tenants", createTenant)
	r.GET("/tenants/:id", getTenant)
	r.POST("/sellers", createSeller)
	r.GET("/sellers/:id", getSeller)

	r.POST("/categories", createCategory)
	r.GET("/categories/:id", getCategory)

	r.POST("/hubs", createHub)
	r.GET("/hubs/:id", getHub)
	r.PUT("/hubs/:id", updateHub)
	r.DELETE("/hubs/:id", deleteHub)
	r.GET("/hubs", listHubs)

	r.POST("/skus", createSKU)
	r.GET("/skus/:id", getSKU)
	r.PUT("/skus/:id", updateSKU)
	r.DELETE("/skus/:id", deleteSKU)
	r.GET("/skus", listSKUs)

	r.PUT("/inventory", upsertInventory)
	r.GET("/inventory", listInventory)

	r.POST("/webhooks", createWebhook)
	r.GET("/webhooks/:id", getWebhook)
	r.PUT("/webhooks/:id", updateWebhook)
	r.DELETE("/webhooks/:id", deleteWebhook)
}
