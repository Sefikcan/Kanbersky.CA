package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sefikcan/kanbersky.ca/internal/server"
	"github.com/sefikcan/kanbersky.ca/pkg/config"
	"github.com/sefikcan/kanbersky.ca/pkg/logger"
	"github.com/sefikcan/kanbersky.ca/pkg/storage/postgres"
	"github.com/sefikcan/kanbersky.ca/pkg/storage/redis"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"log"
)

// @title Go Clean Arch
// @version 1.0
// @description Go Clean Arch
// @contact.name Sefik Can Kanber
// @contact.url https://github.com/sefikcan
// @BasePath /api/v1
// @host localhost:5000
func main()  {
	log.Println("Starting api server")

	cfg := config.NewConfig()

	zapLogger := logger.NewLogger(cfg)
	zapLogger.InitLogger()
	zapLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, false)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	db, err := psqlDB.DB()
	if err != nil {
		zapLogger.Fatalf("Postgresql init: %s", err)
	} else {
		zapLogger.Infof("Postgres connected, Status: %#v", db.Stats())
	}
	defer db.Close()

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	zapLogger.Info("Redis Connected")

	jaegerConfigInstance := jaegerCfg.Configuration{
		ServiceName: cfg.Metric.ServiceName,
		Sampler: &jaegerCfg.SamplerConfig{
			Type: jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans: cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	tracer, closer, err := jaegerConfigInstance.NewTracer(
		jaegerCfg.Logger(jaegerLog.StdLogger),
		jaegerCfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatal("cannot create tracer", err)
	}
	zapLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	zapLogger.Info("Opentracing connected")

	s := server.NewServer(cfg, psqlDB, redisClient, zapLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
