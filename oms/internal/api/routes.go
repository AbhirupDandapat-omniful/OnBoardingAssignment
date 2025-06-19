package api

import "github.com/gin-gonic/gin"

// RegisterRoutes wires up all OMS public endpoints.
func RegisterRoutes(r *gin.Engine) {
    // existing
    r.POST("/orders/bulk", UploadBulkOrders)

    // new: one-step CSV upload endpoint
    r.POST("/orders/upload", UploadCSV)
}
