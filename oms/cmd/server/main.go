package main

import (
	"fmt"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/health"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"

	"github.com/abhirup.dandapat/oms/internal/api"
)

func main() {
	if err := config.Init(30 * time.Second); err != nil {
		panic(err)
	}
	ctx, err := config.TODOContext()
	if err != nil {
		panic(err)
	}

	level := config.GetString(ctx, "log.level")
	log.SetLevel(level)

	logOpts := http.LoggingMiddlewareOptions{
		Format:      config.GetString(ctx, "log.format"),
		Level:       level,
		LogRequest:  true,
		LogResponse: true,
		LogHeader:   false,
	}

	addr := fmt.Sprintf(":%d", config.GetInt(ctx, "server.port"))
	log.Infof("OMS starting, listening on %s", addr)

	srv := http.InitializeServer(
		addr,
		config.GetDuration(ctx, "server.readTimeout"),
		config.GetDuration(ctx, "server.writeTimeout"),
		config.GetDuration(ctx, "server.idleTimeout"),
		false,
		env.RequestID(),
		env.Middleware(config.GetString(ctx, "env")),
		config.Middleware(),
		http.RequestLogMiddleware(logOpts),
	)

	srv.Engine.GET("/health", health.HealthcheckHandler())

	api.RegisterRoutes(srv.Engine)

	if err := srv.StartServer("OMS"); err != nil {
		log.Errorf("OMS shutdown error: %v", err)
	} else {
		log.Infof("OMS stopped gracefully")
	}
}
