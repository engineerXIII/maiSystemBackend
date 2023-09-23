package main

import (
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/server"
	"github.com/engineerXIII/maiSystemBackend/pkg/db/postgres"
	"github.com/engineerXIII/maiSystemBackend/pkg/db/redis"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"log"
	"os"
)

func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("CONFIG_TYPE"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}

	defer psqlDB.Close()

	driver, err := migratePostgres.WithInstance(psqlDB.DB, &migratePostgres.Config{})
	if err != nil {
		appLogger.Fatalf("Cannot create migration driver: %s", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		appLogger.Fatalf("Error on initiate migration: %s", err)
	}
	status := migration.Up()
	if status != nil {
		appLogger.Infof("Migration status: %s", status)
	}
	appLogger.Info("Migration completed")

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	s := server.NewServer(cfg, psqlDB, redisClient, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
