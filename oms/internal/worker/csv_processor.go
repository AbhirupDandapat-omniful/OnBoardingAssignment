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
    awss3  "github.com/aws/aws-sdk-go-v2/service/s3"

    "github.com/omniful/go_commons/config"
    "github.com/omniful/go_commons/log"
    "github.com/omniful/go_commons/sqs"
    commoncsv "github.com/omniful/go_commons/csv"
    gooms3   "github.com/omniful/go_commons/s3"

    "github.com/abhirup.dandapat/oms/constants"
    "github.com/abhirup.dandapat/oms/internal/models"
)

type queueHandler struct{}

func (h *queueHandler) Process(ctx context.Context, msgs *[]sqs.Message) error {
    logger := log.DefaultLogger()

    for _, msg := range *msgs {
        // 1) Unmarshal SQS event
        var evt struct{ Bucket, Key string }
        logger.Infof("raw SQS msg body: %q", string(msg.Value))
        if err := json.Unmarshal(msg.Value, &evt); err != nil {
            logger.Errorf("invalid message JSON: %v", err)
            return err
        }
        logger.Infof("processing S3 object: %s/%s", evt.Bucket, evt.Key)

        s3Client, err := gooms3.NewDefaultAWSS3Client()
        if err != nil {
            logger.Errorf("failed to init S3 client: %v", err)
            return err
        }
        out, err := s3Client.GetObject(ctx, &awss3.GetObjectInput{
            Bucket: awsV2.String(evt.Bucket),
            Key:    awsV2.String(evt.Key),
        })
        if err != nil {
            logger.Errorf("S3 download error: %v", err)
            return err
        }
        defer out.Body.Close()

        r := csv.NewReader(out.Body)
        r.Comma = commoncsv.CsvDelimiter
        r.LazyQuotes = true

        header, err := r.Read()
        if err != nil {
            logger.Errorf("error reading CSV header: %v", err)
            return err
        }
        rows, err := r.ReadAll()
        if err != nil {
            logger.Errorf("error reading CSV rows: %v", err)
            return err
        }

        idx := map[string]int{}
        for i, col := range header {
            idx[col] = i
        }
        for _, col := range []string{"tenant_id", "seller_id", "hub_id", "sku_id", "quantity"} {
            if _, ok := idx[col]; !ok {
                return fmt.Errorf("missing required column %q", col)
            }
        }

        var valid []models.Order
        var invalid [][]string

        for _, row := range rows {
            qstr := row[idx["quantity"]]
            qty, err := strconv.Atoi(qstr)
            if err != nil {
                invalid = append(invalid, row)
                continue
            }
            valid = append(valid, models.Order{
                TenantID: row[idx["tenant_id"]],
                SellerID: row[idx["seller_id"]],
                HubID:    row[idx["hub_id"]],
                SKUID:    row[idx["sku_id"]],
                Quantity: qty,
            })
        }

        for i, o := range valid {
            logger.Infof("valid order %d: %+v", i+1, o)
        }

        if len(invalid) > 0 {
            buf := &bytes.Buffer{}
            w := csv.NewWriter(buf)
            w.Write(header)
            w.WriteAll(invalid)
            w.Flush()

            reader := bytes.NewReader(buf.Bytes())
            errKey := fmt.Sprintf("errors/%s-%d.csv", path.Base(evt.Key), time.Now().Unix())
            if _, err := s3Client.PutObject(ctx, &awss3.PutObjectInput{
                Bucket: awsV2.String(evt.Bucket),
                Key:    awsV2.String(errKey),
                Body:   reader,
            }); err != nil {
                logger.Errorf("failed to upload invalid CSV: %v", err)
            } else {
                logger.Infof("invalid rows CSV at s3://%s/%s", evt.Bucket, errKey)
            }
        }
    }

    return nil
}

func StartCSVProcessor(ctx context.Context) {
    logger := log.DefaultLogger()
    queueURL := config.GetString(ctx, constants.ConfigKeyBulkOrderQueueURL)
    logger.Infof("Listening for SQS messages on %s", queueURL)

    queueName := path.Base(queueURL)
    sqsCfg := &sqs.Config{
        Account:  config.GetString(ctx, "sqs.account"),
        Endpoint: config.GetString(ctx, "sqs.endpoint"),
        Region:   config.GetString(ctx, "aws.region"),
    }
    queueObj, err := sqs.NewStandardQueue(ctx, queueName, sqsCfg)
    if err != nil {
        logger.Panicf("failed to create SQS queue object: %v", err)
    }
    consumer, err := sqs.NewConsumer(queueObj, 1, 1, &queueHandler{}, 10, 30, false, false)
    if err != nil {
        logger.Panicf("failed to init SQS consumer: %v", err)
    }
    consumer.Start(ctx)
}
