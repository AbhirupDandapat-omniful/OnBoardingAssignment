package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/abhirup.dandapat/oms/internal/models"
)

func collection(c *gin.Context) (*mongo.Collection, error) {
	ctx := c.Request.Context()
	uri := config.GetString(ctx, "mongo.uri")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client.Database("omsdb").Collection("webhooks"), nil
}

func createWebhook(c *gin.Context) {
	var w models.Webhook
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	now := time.Now().UTC()
	w.ID = uuid.New().String()
	w.CreatedAt, w.UpdatedAt = now, now
	w.IsActive = true

	coll, err := collection(c)
	if err != nil {
		log.DefaultLogger().Errorf("createWebhook: connect db: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	if _, err := coll.InsertOne(c.Request.Context(), w); err != nil {
		log.DefaultLogger().Errorf("createWebhook: insert: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.create_webhook_failed")})
		return
	}
	c.JSON(http.StatusCreated, w)
}

func getWebhook(c *gin.Context) {
	id := c.Param("id")
	coll, err := collection(c)
	if err != nil {
		log.DefaultLogger().Errorf("getWebhook: connect db: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	var w models.Webhook
	if err := coll.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&w); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": i18n.Translate(c, "error.webhook_not_found")})
		return
	}
	c.JSON(http.StatusOK, w)
}

func listWebhooks(c *gin.Context) {
	tenant := c.Query("tenant_id")
	filter := bson.M{"is_active": true}
	if tenant != "" {
		filter["tenant_id"] = tenant
	}

	coll, err := collection(c)
	if err != nil {
		log.DefaultLogger().Errorf("listWebhooks: connect db: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	cursor, err := coll.Find(c.Request.Context(), filter)
	if err != nil {
		log.DefaultLogger().Errorf("listWebhooks: find: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	defer cursor.Close(c.Request.Context())

	var all []models.Webhook
	for cursor.Next(c.Request.Context()) {
		var w models.Webhook
		if err := cursor.Decode(&w); err != nil {
			log.DefaultLogger().Errorf("listWebhooks: decode: %v", err)
			continue
		}
		all = append(all, w)
	}
	c.JSON(http.StatusOK, all)
}

func updateWebhook(c *gin.Context) {
	id := c.Param("id")
	var update models.Webhook
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n.Translate(c, "error.invalid_request")})
		return
	}
	update.UpdatedAt = time.Now().UTC()

	coll, err := collection(c)
	if err != nil {
		log.DefaultLogger().Errorf("updateWebhook: connect db: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	_, err = coll.UpdateOne(
		c.Request.Context(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"callback_url": update.CallbackURL,
			"events":       update.Events,
			"headers":      update.Headers,
			"is_active":    update.IsActive,
			"updated_at":   update.UpdatedAt,
		}},
	)
	if err != nil {
		log.DefaultLogger().Errorf("updateWebhook: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.update_webhook_failed")})
		return
	}
	c.Status(http.StatusOK)
}

func deleteWebhook(c *gin.Context) {
	id := c.Param("id")
	coll, err := collection(c)
	if err != nil {
		log.DefaultLogger().Errorf("deleteWebhook: connect db: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.internal")})
		return
	}
	// soft‐delete by marking inactive
	_, err = coll.UpdateOne(
		c.Request.Context(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_active": false, "updated_at": time.Now().UTC()}},
	)
	if err != nil {
		log.DefaultLogger().Errorf("deleteWebhook: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n.Translate(c, "error.delete_webhook_failed")})
		return
	}
	c.Status(http.StatusNoContent)
}
