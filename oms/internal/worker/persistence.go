package worker

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/abhirup.dandapat/oms/internal/models"
)

var persistLogger = log.DefaultLogger()

func getMongoClient(ctx context.Context) (*mongo.Client, error) {
	uri := config.GetString(ctx, "mongo.uri")
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

func getOrdersCollection(ctx context.Context) (*mongo.Collection, error) {
	client, err := getMongoClient(ctx)
	if err != nil {
		return nil, err
	}
	uri := config.GetString(ctx, "mongo.uri")
	parts := strings.Split(strings.TrimSuffix(uri, "/"), "/")
	dbName := parts[len(parts)-1]
	return client.Database(dbName).Collection("orders"), nil
}

func saveOrder(ctx context.Context, o *models.Order) error {
	coll, err := getOrdersCollection(ctx)
	if err != nil {
		persistLogger.Errorf("Mongo init error: %v", err)
		return err
	}
	o.ID = uuid.NewString()
	o.Status = "on_hold"
	o.CreatedAt = time.Now().UTC()
	if _, err := coll.InsertOne(ctx, o); err != nil {
		persistLogger.Errorf("Mongo InsertOne error: %v", err)
		return err
	}
	return nil
}
