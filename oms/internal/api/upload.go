package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	gooms3 "github.com/omniful/go_commons/s3"
	"github.com/omniful/go_commons/sqs"

	// "github.com/abhirup.dandapat/oms/constants"
	"github.com/abhirup.dandapat/oms/internal/store"
)

func UploadCSV(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	defer file.Close()

	bucket := config.GetString(c.Request.Context(), "s3.uploadBucket")
	if bucket == "" {
		log.DefaultLogger().Error("UploadCSV: upload bucket not configured")
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	key := fmt.Sprintf("uploads/%s-%s", time.Now().Format("20060102-150405"),
		uuid.New().String()+"_"+header.Filename)

	s3Client, err := gooms3.NewDefaultAWSS3Client()
	if err != nil {
		log.DefaultLogger().Errorf("UploadCSV: init S3 client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}

	putIn := &awss3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	}
	if _, err := s3Client.PutObject(c.Request.Context(), putIn); err != nil {
		log.DefaultLogger().Errorf("UploadCSV: PutObject error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}

	publisher, err := store.NewBulkOrderPublisher(c.Request.Context())
	if err != nil {
		log.DefaultLogger().Errorf("UploadCSV: init publisher: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	evt := struct {
		Bucket string `json:"bucket"`
		Key    string `json:"key"`
	}{Bucket: bucket, Key: key}
	payload, _ := json.Marshal(evt)
	msg := &sqs.Message{Value: payload}
	if err := publisher.Publish(c.Request.Context(), msg); err != nil {
		log.DefaultLogger().Errorf("UploadCSV: Publish error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.enqueue_bulk_order_failed")})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"bucket": bucket, "key": key})
}
