// internal/api/errors.go
package api

import (
	"fmt"
	"net/http"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	gooms3 "github.com/omniful/go_commons/s3"
)

var errorLogger = log.DefaultLogger()

func DownloadErrorCSV(c *gin.Context) {
	file := c.Param("file")
	bucket := config.GetString(c.Request.Context(), "s3.uploadBucket")
	key := fmt.Sprintf("errors/%s", file)

	s3Client, err := gooms3.NewDefaultAWSS3Client()
	if err != nil {
		errorLogger.Errorf("DownloadErrorCSV: S3 init error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	out, err := s3Client.GetObject(c.Request.Context(), &awss3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		errorLogger.Errorf("DownloadErrorCSV: GetObject error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer out.Body.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
	length := int64(0)
	if out.ContentLength != nil {
		length = *out.ContentLength
	}
	c.DataFromReader(http.StatusOK, length, "text/csv", out.Body, nil)
}
