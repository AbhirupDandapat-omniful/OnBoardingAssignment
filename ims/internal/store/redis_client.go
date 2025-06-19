package store

import (
    "context"
    "time"

    "github.com/omniful/go_commons/config"
    "github.com/omniful/go_commons/log"
    "github.com/omniful/go_commons/redis"
)

var RedisClient *redis.Client

func InitRedis(ctx context.Context) {
    logger := log.DefaultLogger()

    endpoint := config.GetString(ctx, "redis.endpoint")
    dbIndex := config.GetInt(ctx, "redis.db")

    cfg := &redis.Config{
        Hosts: []string{endpoint},
        DB:    uint(dbIndex),
    }

    RedisClient = redis.NewClient(cfg)

    pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    if err := RedisClient.Ping(pingCtx).Err(); err != nil {
        logger.Panicf("Redis ping failed: %v", err)
    }

    logger.Infof("Connected to Redis at %s (db=%d)", endpoint, dbIndex)
}
