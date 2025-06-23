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

func validateEntity(ctx context.Context, cache *redis.Client, client *commonsHttp.Client, baseURL, prefix, id string) bool {
	cacheKey := prefix + ":" + id

	if _, err := cache.Get(ctx, cacheKey); err == nil {
		return true
	}

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

	cache.Set(ctx, cacheKey, "1", 5*time.Minute)
	return true
}

func setupValidationClients(ctx context.Context) (*redis.Client, *commonsHttp.Client, string, error) {
	addrList := config.GetStringSlice(ctx, "redis.addrs")
	redisCfg := &redis.Config{Hosts: addrList}
	cache := redis.NewClient(redisCfg)

	baseURL := config.GetString(ctx, "ims.baseUrl")

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
