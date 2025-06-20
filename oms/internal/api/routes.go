package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	// existing:
	r.POST("/orders/bulk", UploadBulkOrders)
	r.POST("/orders/upload", UploadCSV)

	// new:
	r.GET("/orders/errors/:file", DownloadErrorCSV)
}
