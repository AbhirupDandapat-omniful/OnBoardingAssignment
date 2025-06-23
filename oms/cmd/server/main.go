package main

import (
	"strconv"
	"time"

	"github.com/abhirup.dandapat/oms/internal/api"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/health"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
)

func main() {
	if err := config.Init(30 * time.Second); err != nil {
		panic(err)
	}
	ctx, err := config.TODOContext()
	if err != nil {
		panic(err)
	}

	lvl := config.GetString(ctx, "log.level")
	log.SetLevel(lvl)
	log.Infof("Starting OMS on port %d", config.GetInt(ctx, "server.port"))

	port := ":" + strconv.Itoa(config.GetInt(ctx, "server.port"))
	srv := http.InitializeServer(
		port,
		config.GetDuration(ctx, "server.readTimeout"),
		config.GetDuration(ctx, "server.writeTimeout"),
		config.GetDuration(ctx, "server.idleTimeout"),
		false,
		env.RequestID(),
		env.Middleware(config.GetString(ctx, "env")),
	)

	srv.Engine.GET("/health", health.HealthcheckHandler())
	api.RegisterRoutes(srv.Engine)

	if err := srv.StartServer("OMS"); err != nil {
		log.Errorf("OMS shutdown error: %v", err)
	}
}
