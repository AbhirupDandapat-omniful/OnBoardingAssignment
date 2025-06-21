package worker

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"time"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/omniful/go_commons/config"
	commoncsv "github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/log"
	gooms3 "github.com/omniful/go_commons/s3"
	"github.com/omniful/go_commons/sqs"

	"github.com/abhirup.dandapat/oms/constants"
	"github.com/abhirup.dandapat/oms/internal/models"
)

type queueHandler struct{}

func (h *queueHandler) Process(ctx context.Context, msgs *[]sqs.Message) error {
	logger := log.DefaultLogger()

	s3Client, err := gooms3.NewDefaultAWSS3Client()
	if err != nil {
		logger.Errorf("failed to init S3 client: %v", err)
		return err
	}

	producer := newProducer(ctx)

	for _, msg := range *msgs {
		var evt struct{ Bucket, Key string }
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			logger.Errorf("invalid SQS JSON: %v", err)
			return err
		}
		logger.Infof("processing S3 object: %s/%s", evt.Bucket, evt.Key)

		out, err := s3Client.GetObject(ctx, &awss3.GetObjectInput{
			Bucket: awsV2.String(evt.Bucket),
			Key:    awsV2.String(evt.Key),
		})
		if err != nil {
			logger.Errorf("S3 GetObject error: %v", err)
			return err
		}
		defer out.Body.Close()

		r := csv.NewReader(out.Body)
		r.Comma = commoncsv.CsvDelimiter
		r.LazyQuotes = true

		header, err := r.Read()
		if err != nil {
			logger.Errorf("read header: %v", err)
			return err
		}
		rows, err := r.ReadAll()
		if err != nil {
			logger.Errorf("read rows: %v", err)
			return err
		}

		idx := map[string]int{}
		for i, col := range header {
			idx[col] = i
		}

		var invalid [][]string

		for _, row := range rows {
			qty, err := strconv.Atoi(row[idx["quantity"]])
			if err != nil || qty <= 0 {
				invalid = append(invalid, row)
				continue
			}

			order := &models.Order{
				TenantID: row[idx["tenant_id"]],
				SellerID: row[idx["seller_id"]],
				HubID:    row[idx["hub_id"]],
				SKUID:    row[idx["sku_id"]],
				Quantity: int64(qty),
			}

			if err := saveOrder(ctx, order); err != nil {
				invalid = append(invalid, row)
				continue
			}

			publishOrderCreated(ctx, producer, order)
			logger.Infof("Processed order: %+v", order)
		}

		if len(invalid) > 0 {
			buf := &bytes.Buffer{}
			w := csv.NewWriter(buf)
			w.Write(header)
			w.WriteAll(invalid)
			w.Flush()

			errKey := fmt.Sprintf("errors/%s-%d.csv",
				path.Base(evt.Key),
				time.Now().Unix(),
			)
			if _, err := s3Client.PutObject(ctx, &awss3.PutObjectInput{
				Bucket: awsV2.String(evt.Bucket),
				Key:    awsV2.String(errKey),
				Body:   bytes.NewReader(buf.Bytes()),
			}); err != nil {
				logger.Errorf("upload invalid CSV: %v", err)
			} else {
				logger.Infof("invalid rows CSV: s3://%s/%s", evt.Bucket, errKey)
			}
		}
	}

	return nil
}

func StartCSVProcessor(ctx context.Context) {
	log.DefaultLogger().Infof("CONFIG kafka.brokers = %#v", config.GetStringSlice(ctx, "kafka.brokers"))
	log.DefaultLogger().Infof("CONFIG kafka.version = %q", config.GetString(ctx, "kafka.version"))

	logger := log.DefaultLogger()
	queueURL := config.GetString(ctx, constants.ConfigKeyBulkOrderQueueURL)
	logger.Infof("Listening for SQS on %s", queueURL)

	qName := path.Base(queueURL)
	sqsCfg := &sqs.Config{
		Account:  config.GetString(ctx, "sqs.account"),
		Endpoint: config.GetString(ctx, "sqs.endpoint"),
		Region:   config.GetString(ctx, "aws.region"),
	}
	qObj, err := sqs.NewStandardQueue(ctx, qName, sqsCfg)
	if err != nil {
		logger.Panicf("failed to create SQS queue: %v", err)
	}
	consumer, err := sqs.NewConsumer(
		qObj,
		uint64(config.GetInt(ctx, "sqs.consumer.workerCount")),
		uint64(config.GetInt(ctx, "sqs.consumer.concurrencyPerWorker")),
		&queueHandler{},
		int64(config.GetInt(ctx, "sqs.consumer.batchSize")),
		int64(config.GetInt(ctx, "sqs.consumer.visibilityTimeout")),
		false,
		false,
	)
	if err != nil {
		logger.Panicf("failed to init SQS consumer: %v", err)
	}
	consumer.Start(ctx)
}
