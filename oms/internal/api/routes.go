package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	r.POST("/orders/bulk", UploadBulkOrders)
	r.POST("/orders/upload", UploadCSV)
	r.GET("/orders/errors/:file", DownloadErrorCSV)

	// New: filtered order-list endpoint
	r.GET("/orders", ListOrders)
}
