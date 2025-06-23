package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	// CSV + bulk
	r.POST("/orders/bulk", UploadBulkOrders)
	r.POST("/orders/upload", UploadCSV)

	// new public endpoints
	r.GET("/orders", ListOrders)
	r.GET("/orders/errors/:file", DownloadErrorCSV)
	r.POST("/orders", CreateOrder)
	// your existing webhook CRUD, etc...
	r.POST("/webhooks", createWebhook)
	r.GET("/webhooks/:id", getWebhook)
	r.GET("/webhooks", listWebhooks)
	r.PUT("/webhooks/:id", updateWebhook)
	r.DELETE("/webhooks/:id", deleteWebhook)

}
