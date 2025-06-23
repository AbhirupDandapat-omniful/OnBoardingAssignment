package worker

import (
	"context"

	"github.com/abhirup.dandapat/oms/internal/models"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
)

var kafkaLogger = log.DefaultLogger()

func newProducer(ctx context.Context) *kafka.ProducerClient {
	brokers := config.GetStringSlice(ctx, "kafka.brokers")
	clientID := config.GetString(ctx, "kafka.clientId")
	version := config.GetString(ctx, "kafka.version")
	return kafka.NewProducer(
		kafka.WithBrokers(brokers),
		kafka.WithClientID(clientID),
		kafka.WithKafkaVersion(version),
	)
}

func publishOrderCreated(ctx context.Context, producer *kafka.ProducerClient, o *models.Order) {
	payload, err := pubsub.NewEventInBytes(models.OrderCreated{
		OrderID:   o.ID,
		TenantID:  o.TenantID,
		SellerID:  o.SellerID,
		HubID:     o.HubID,
		SKUID:     o.SKUID,
		Quantity:  o.Quantity,
		CreatedAt: o.CreatedAt,
	})
	if err != nil {
		kafkaLogger.Errorf("marshal OrderCreated: %v", err)
		return
	}
	msg := &pubsub.Message{
		Topic: config.GetString(ctx, "kafka.topicOrderCreated"),
		Key:   o.ID,
		Value: payload,
	}
	if err := producer.Publish(ctx, msg); err != nil {
		kafkaLogger.Errorf("Kafka publish error: %v", err)
	}
}
