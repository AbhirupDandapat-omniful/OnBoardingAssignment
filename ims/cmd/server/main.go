package main

import (
    "fmt"
    "time"

    "github.com/omniful/go_commons/config"
    "github.com/omniful/go_commons/env"
    "github.com/omniful/go_commons/health"
    "github.com/omniful/go_commons/http"
    "github.com/omniful/go_commons/log"

    "github.com/abhirup.dandapat/ims/internal/api"
    "github.com/abhirup.dandapat/ims/internal/store"
)

func main() {
    // 1) Load config (hot-reload every 30s)
    if err := config.Init(30 * time.Second); err != nil {
        panic(err)
    }
    ctx, err := config.TODOContext()
    if err != nil {
        panic(err)
    }

    // 2) Init Postgres + migrations
    store.InitPostgres(ctx)

    // 3) Init Redis
    store.InitRedis(ctx)

    // 4) Configure logging
    log.SetLevel(config.GetString(ctx, "log.level"))
    port := config.GetInt(ctx, "server.port")
    log.Infof("Starting IMS on port %d", port)

    // 5) HTTP server + middleware: request-ID, env, i18n, etc.
    srv := http.InitializeServer(
        fmt.Sprintf(":%d", port),
        config.GetDuration(ctx, "server.readTimeout"),
        config.GetDuration(ctx, "server.writeTimeout"),
        config.GetDuration(ctx, "server.idleTimeout"),
        false,
        env.RequestID(),
        env.Middleware(config.GetString(ctx, "env")),
        config.Middleware(), // i18n + error wrapping
    )

    // 6) Health endpoint
    srv.Engine.GET("/health", health.HealthcheckHandler())

    // 7) Register all business routes
    api.RegisterRoutes(srv.Engine)

    // 8) Run (blocking)
    if err := srv.StartServer("IMS"); err != nil {
        log.Errorf("IMS shutdown error: %v", err)
    }
}
