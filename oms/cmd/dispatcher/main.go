package main

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/omniful/go_commons/config"
	commonsHttp "github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/abhirup.dandapat/oms/internal/dispatcher"
	"github.com/abhirup.dandapat/oms/internal/webhooklogs"
)

func main() {

	if err := config.Init(10 * time.Second); err != nil {
		log.DefaultLogger().Panicf("config init failed: %v", err)
	}
	ctx, err := config.TODOContext()
	if err != nil {
		log.DefaultLogger().Panicf("getting config context failed: %v", err)
	}

	log.DefaultLogger().Info("Starting Webhook Dispatcher…")

	mongoURI := config.GetString(ctx, "mongo.uri")
	mcli, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.DefaultLogger().Panicf("mongo connect failed: %v", err)
	}
	coll := mcli.Database("omsdb").Collection("webhooks")

	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient, err := commonsHttp.NewHTTPClient(
		config.GetString(ctx, "kafka.clientId"),
		"", transport,
	)
	if err != nil {
		log.DefaultLogger().Panicf("http client init failed: %v", err)
	}

	producer := kafka.NewProducer(
		kafka.WithBrokers(config.GetStringSlice(ctx, "kafka.brokers")),
		kafka.WithClientID(config.GetString(ctx, "kafka.clientId")+"-producer"),
		kafka.WithKafkaVersion(config.GetString(ctx, "kafka.version")),
	)

	consumer := kafka.NewConsumer(
		kafka.WithBrokers(config.GetStringSlice(ctx, "kafka.brokers")),
		kafka.WithConsumerGroup(config.GetString(ctx, "kafka.groupId")),
		kafka.WithClientID(config.GetString(ctx, "kafka.clientId")),
		kafka.WithKafkaVersion(config.GetString(ctx, "kafka.version")),
	)

	createdTopic := config.GetString(ctx, "kafka.topicOrderCreated")
	updatedTopic := config.GetString(ctx, "kafka.topicOrderUpdated")
	failedTopic := config.GetString(ctx, "kafka.topicWebhookFailed")

	rawDisp := dispatcher.NewDispatcher(coll, httpClient, producer, failedTopic)
	retryDisp := dispatcher.NewRetryHandler(rawDisp, 3, time.Second)

	consumer.RegisterHandler(createdTopic, retryDisp)
	consumer.RegisterHandler(updatedTopic, retryDisp)

	logColl := mcli.Database("omsdb").Collection("webhook_logs")
	logHandler := webhooklogs.NewHandler(logColl)
	retryLogHandler := dispatcher.NewRetryHandler(logHandler, 3, time.Second)
	consumer.RegisterHandler(config.GetString(ctx, "kafka.topicWebhookFailed"), retryLogHandler)

	consumer.Subscribe(ctx)

	select {}
}
