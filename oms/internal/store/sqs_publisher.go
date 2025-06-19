package store

import (
    "context"
    "path"

    "github.com/omniful/go_commons/config"
    "github.com/omniful/go_commons/log"
    "github.com/omniful/go_commons/sqs"

    "github.com/abhirup.dandapat/oms/constants"
)

// NewBulkOrderPublisher returns a GoCommons SQS publisher for the bulk-order queue.
func NewBulkOrderPublisher(ctx context.Context) (*sqs.Publisher, error) {
    queueURL := config.GetString(ctx, constants.ConfigKeyBulkOrderQueueURL)
    queueName := path.Base(queueURL)

    sqsCfg := &sqs.Config{
        Account:  config.GetString(ctx, "sqs.account"),
        Endpoint: config.GetString(ctx, "sqs.endpoint"),
        Region:   config.GetString(ctx, "aws.region"),
    }
    q, err := sqs.NewStandardQueue(ctx, queueName, sqsCfg)
    if err != nil {
        log.DefaultLogger().Errorf("NewBulkOrderPublisher: NewStandardQueue error: %v", err)
        return nil, err
    }
    return sqs.NewPublisher(q), nil
}
