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

	"github.com/abhirup.dandapat/oms/internal/finalizer"
)

func main() {

	if err := config.Init(10 * time.Second); err != nil {
		log.DefaultLogger().Panicf("config init failed: %v", err)
	}
	ctx, err := config.TODOContext()
	if err != nil {
		log.DefaultLogger().Panicf("getting config context failed: %v", err)
	}

	log.DefaultLogger().Info("Starting Order Finalizer…")

	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	baseURL := config.GetString(ctx, "ims.baseUrl")
	httpClient, err := commonsHttp.NewHTTPClient(
		config.GetString(ctx, "finalizer.clientID"),
		baseURL,
		transport,
	)
	if err != nil {
		log.DefaultLogger().Panicf("http client init failed: %v", err)
	}

	mongoURI := config.GetString(ctx, "mongo.uri")
	mcli, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.DefaultLogger().Panicf("mongo connect failed: %v", err)
	}
	ordersColl := mcli.Database("omsdb").Collection("orders")

	producer := kafka.NewProducer(
		kafka.WithBrokers(config.GetStringSlice(ctx, "kafka.brokers")),
		kafka.WithClientID(config.GetString(ctx, "kafka.clientId")+"-producer"),
		kafka.WithKafkaVersion(config.GetString(ctx, "kafka.version")),
	)

	rawH := finalizer.NewHandler(ordersColl, httpClient, producer)
	retryH := finalizer.NewRetryHandler(rawH, 3, 1*time.Second)

	consumer := kafka.NewConsumer(
		kafka.WithBrokers(config.GetStringSlice(ctx, "kafka.brokers")),
		kafka.WithConsumerGroup(config.GetString(ctx, "finalizer.groupID")),
		kafka.WithClientID(config.GetString(ctx, "finalizer.clientID")),
		kafka.WithKafkaVersion(config.GetString(ctx, "kafka.version")),
	)
	topic := config.GetString(ctx, "finalizer.topicCreated")
	consumer.RegisterHandler(topic, retryH)

	consumer.Subscribe(ctx)

	select {}
}
