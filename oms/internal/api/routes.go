package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {

	r.POST("/orders/bulk", UploadBulkOrders)
	r.POST("/orders/upload", UploadCSV)

	r.GET("/orders", ListOrders)
	r.GET("/orders/errors/:file", DownloadErrorCSV)
	r.POST("/orders", CreateOrder)
	r.POST("/webhooks", createWebhook)
	r.GET("/webhooks/:id", getWebhook)
	r.GET("/webhooks", listWebhooks)
	r.PUT("/webhooks/:id", updateWebhook)
	r.DELETE("/webhooks/:id", deleteWebhook)

	r.GET("/webhook-logs", ListWebhookLogs)

}
