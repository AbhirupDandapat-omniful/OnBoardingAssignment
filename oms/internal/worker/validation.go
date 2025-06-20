package worker

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"time"

	"github.com/omniful/go_commons/config"
	commonsHttp "github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/redis"
)

var validateLogger = log.DefaultLogger()

// validateEntity checks Redis cache, then calls IMS at GET {baseURL}/{prefix}/{id}.
func validateEntity(ctx context.Context, cache *redis.Client, client *commonsHttp.Client, baseURL, prefix, id string) bool {
	cacheKey := prefix + ":" + id

	// 1) Redis cache
	if _, err := cache.Get(ctx, cacheKey); err == nil {
		return true
	}

	// 2) IMS HTTP call
	req := &commonsHttp.Request{
		Url:     fmt.Sprintf("/%s/%s", prefix, id),
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(req, nil)
	if err != nil {
		validateLogger.Warnf("IMS GET %s failed: %v", req.Url, err)
		return false
	}
	if resp.StatusCode() != stdhttp.StatusOK {
		validateLogger.Warnf("IMS GET %s returned %d", req.Url, resp.StatusCode())
		return false
	}

	// 3) Cache hit for 5 minutes
	cache.Set(ctx, cacheKey, "1", 5*time.Minute)
	return true
}

// setupValidationClients returns a Redis client, an HTTP client for IMS, and the IMS base URL.
func setupValidationClients(ctx context.Context) (*redis.Client, *commonsHttp.Client, string, error) {
	// Redis
	addrList := config.GetStringSlice(ctx, "redis.addrs")
	redisCfg := &redis.Config{Hosts: addrList}
	cache := redis.NewClient(redisCfg)

	// IMS base URL
	baseURL := config.GetString(ctx, "ims.baseUrl")

	// HTTP client
	httpClient, err := commonsHttp.NewHTTPClient(
		"oms-validation",
		baseURL,
		nil,
		commonsHttp.WithTimeout(10*time.Second),
	)
	if err != nil {
		return nil, nil, "", err
	}

	return cache, httpClient, baseURL, nil
}
