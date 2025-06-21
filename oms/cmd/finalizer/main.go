// cmd/finalizer/main.go
package main

import (
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"

	"github.com/abhirup.dandapat/oms/internal/finalizer"
)

func main() {
	// 1) Load config (poll every 10s for changes)
	if err := config.Init(10 * time.Second); err != nil {
		log.DefaultLogger().Panicf("config init failed: %v", err)
	}

	// 2) Get a context that carries your loaded config
	ctx, err := config.TODOContext()
	if err != nil {
		log.DefaultLogger().Panicf("getting config context failed: %v", err)
	}

	log.DefaultLogger().Info("Starting Order Finalizer…")

	// 3) Build the Kafka consumer
	consumer := kafka.NewConsumer(
		kafka.WithBrokers(config.GetStringSlice(ctx, "kafka.brokers")),
		kafka.WithConsumerGroup(config.GetString(ctx, "kafka.groupId")),
		kafka.WithClientID(config.GetString(ctx, "kafka.clientId")),
		kafka.WithKafkaVersion(config.GetString(ctx, "kafka.version")),
	)

	// 4) Register your handler (must implement pubsub.IPubSubMessageHandler)
	handler := finalizer.NewHandler()
	topic := config.GetString(ctx, "kafka.topicOrderCreated")
	consumer.RegisterHandler(topic, handler)

	// 5) Start consuming
	consumer.Subscribe(ctx)

	// 6) Block forever
	select {}
}
