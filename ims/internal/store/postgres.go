package store

import (
    "context"
    "os"
    "path/filepath"
    "time"

    "github.com/omniful/go_commons/config"
    "github.com/omniful/go_commons/log"
    "github.com/omniful/go_commons/db/sql/migration"
    "github.com/omniful/go_commons/db/sql/postgres"
)

var DB *postgres.DbCluster

func InitPostgres(ctx context.Context) {
    logger := log.DefaultLogger()

    master := postgres.DBConfig{
        Host:                   config.GetString(ctx, "postgres.primary.host"),
        Port:                   config.GetString(ctx, "postgres.primary.port"),
        Username:               config.GetString(ctx, "postgres.primary.username"),
        Password:               config.GetString(ctx, "postgres.primary.password"),
        Dbname:                 config.GetString(ctx, "postgres.primary.database"),
        MaxOpenConnections:     config.GetInt(ctx, "postgres.pool.max_open_conns"),
        MaxIdleConnections:     config.GetInt(ctx, "postgres.pool.max_idle_conns"),
        ConnMaxLifetime:        config.GetDuration(ctx, "postgres.pool.conn_max_lifetime"),
        DebugMode:              config.GetBool(ctx, "postgres.primary.debug_mode"),
        PrepareStmt:            config.GetBool(ctx, "postgres.primary.prepare_stmt"),
        SkipDefaultTransaction: config.GetBool(ctx, "postgres.primary.skip_default_transaction"),
    }

    var replicas []postgres.DBConfig

    DB = postgres.InitializeDBInstance(master, &replicas)
    logger.Infof("Postgres cluster initialized: master=%s:%s; replicas=%d",
        master.Host, master.Port, len(replicas),
    )

    wd, err := os.Getwd()
    if err != nil {
        logger.Panicf("cannot get wd: %v", err)
    }
    migrationsDir := filepath.ToSlash(filepath.Join(wd, "sql", "migration"))
    fileURL := "file://" + migrationsDir

    dbURL := migration.BuildSQLDBURL(
        master.Host, master.Port, master.Dbname,
        master.Username, master.Password,
    )

    logger.Infof("Applying migrations from %s", fileURL)
    migration.Execute(fileURL, dbURL, "up", "")

    sqlDB, err := DB.GetMasterDB(ctx).DB()
    if err != nil {
        logger.Panicf("extract sql.DB failed: %v", err)
    }
    pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    if err := sqlDB.PingContext(pingCtx); err != nil {
        logger.Panicf("Postgres ping failed: %v", err)
    }
    logger.Infof("Postgres master ping successful")
}
