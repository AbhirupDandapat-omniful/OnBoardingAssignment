package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	gooms3 "github.com/omniful/go_commons/s3"
	"github.com/omniful/go_commons/sqs"

	"github.com/abhirup.dandapat/oms/internal/store"
)

type BulkOrderRequest struct {
	Bucket string `json:"bucket" binding:"required"`
	Key    string `json:"key"    binding:"required"`
}

type CreateBulkOrderEvent struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

func UploadBulkOrders(c *gin.Context) {
	var req BulkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}

	if strings.Contains(req.Key, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_s3_key")})
		return
	}

	s3Client, err := gooms3.NewDefaultAWSS3Client()
	if err != nil {
		log.DefaultLogger().Errorf("UploadBulkOrders: init S3 client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	if _, err := s3Client.HeadObject(c.Request.Context(), &awss3.HeadObjectInput{
		Bucket: aws.String(req.Bucket),
		Key:    aws.String(req.Key),
	}); err != nil {
		log.DefaultLogger().Warnf("S3 object not found: %s/%s: %v", req.Bucket, req.Key, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("S3 object not found: %s/%s", req.Bucket, req.Key)})
		return
	}

	publisher, err := store.NewBulkOrderPublisher(c.Request.Context())
	if err != nil {
		log.DefaultLogger().Errorf("UploadBulkOrders: init publisher: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	evt := CreateBulkOrderEvent{Bucket: req.Bucket, Key: req.Key}
	payload, err := json.Marshal(evt)
	if err != nil {
		log.DefaultLogger().Errorf("UploadBulkOrders: marshal event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}

	msg := &sqs.Message{Value: payload}
	if err := publisher.Publish(c.Request.Context(), msg); err != nil {
		log.DefaultLogger().Errorf("UploadBulkOrders: Publish error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.enqueue_bulk_order_failed")})
		return
	}

	c.Status(http.StatusAccepted)
}
