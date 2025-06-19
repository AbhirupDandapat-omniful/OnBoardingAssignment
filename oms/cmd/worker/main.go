package main

import (
    "time"

    "github.com/omniful/go_commons/config"
    "github.com/omniful/go_commons/log"
    "github.com/abhirup.dandapat/oms/internal/worker"
)

func main() {
    if err := config.Init(30 * time.Second); err != nil {
        panic(err)
    }
    ctx, err := config.TODOContext()
    if err != nil {
        panic(err)
    }

    log.SetLevel(config.GetString(ctx, "log.level"))
    log.Infof("Starting CSV‐Processor worker")

    worker.StartCSVProcessor(ctx)
    select {} // block forever
}
