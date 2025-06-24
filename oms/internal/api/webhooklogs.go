package api

import (
	"net/http"

	"github.com/abhirup.dandapat/oms/internal/webhooklogs"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListWebhookLogs(c *gin.Context) {
	ctx := c.Request.Context()
	mongoURI := config.GetString(ctx, "mongo.uri")
	clientOpts := options.Client().ApplyURI(mongoURI)
	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.DefaultLogger().Errorf("ListOrders: mongo connect error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	defer mongoClient.Disconnect(ctx)

	coll := mongoClient.Database("omsdb").Collection("webhook_logs")

	cursor, err := coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		c.JSON(500, gin.H{"error": "internal"})
		return
	}
	var logs []webhooklogs.FailedWebhook
	if err := cursor.All(ctx, &logs); err != nil {
		c.JSON(500, gin.H{"error": "internal"})
		return
	}
	c.JSON(200, gin.H{"logs": logs})
}
