package main

import (
	"github.com/engineerXIII/maiSystemBackend/config"
	server "github.com/engineerXIII/maiSystemBackend/internal/service/notification"
	"github.com/engineerXIII/maiSystemBackend/pkg/amqp/rabbitmq"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
	"github.com/go-co-op/gocron"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"log"
	"os"
	"time"
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

	amqpClient, err := rabbitmq.NewAMQP(cfg)
	if err != nil {
		appLogger.Fatalf("Cannot connect to AMQP: %s", err)
	}
	defer func() {
		_ = amqpClient.Close()
	}()
	appLogger.Info("AMQP connected")

	amqpChannel, err := rabbitmq.CreateChannel(amqpClient)
	if err != nil {
		appLogger.Fatalf("Cannot open channel AMQP: %s", err)
	}
	defer func() {
		_ = amqpChannel.Close() // Закрываем канал в случае удачной попытки открытия
	}()
	appLogger.Info("AMQP channel opened")
	amqpQueue, err := rabbitmq.DeclareQueue(amqpChannel, cfg)
	if err != nil {
		appLogger.Fatalf("Cannot open channel AMQP: %s", err)
	}

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

	cron := gocron.NewScheduler(time.UTC)
	appLogger.Info("Cron started")

	s := server.NewServer(cfg, amqpChannel, amqpQueue, cron, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
