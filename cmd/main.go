package main

import (
	"context"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"proxy/config"
	"proxy/internal/infrastructure/db/postgres"
	"proxy/run"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func init() {
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Fatalf("не удалось создать директорию для логов: %v", err)
	}
}

func main() {
	logger := setupLogger()
	defer logger.Sync()

	cfg := config.MustLoad()
	logger.Infow("Конфигурация загружена")

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatalw("Не удалось подключиться к базе данных", "err", err)
	}

	application := run.NewApp(logger, db, cfg)

	tp := initTracer()
	defer func() { _ = tp.Shutdown(context.Background()) }()
	otel.SetTracerProvider(tp)

	// Инициализация реестра метрик
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"go_project": cfg.Log.Project}, registry)
	registerer.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	// Обработка метрик
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	logger.Infow("Завершение работы", "signal", sign)

	application.Stop()
	logger.Info("Приложение остановлено")
}

func setupLogger() *zap.SugaredLogger {
	logConfig := zap.NewProductionConfig()
	logConfig.OutputPaths = []string{
		"stdout",
		"logs/app.log",
	}
	logConfig.ErrorOutputPaths = []string{
		"stderr",
		"logs/error.log",
	}
	logConfig.EncoderConfig.TimeKey = "timestamp"
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatalf("не удалось инициализировать логгер: %v", err)
	}

	return logger.Sugar()
}

func initTracer() *trace.TracerProvider {
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("не удалось создать экспортер: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("proxy-service"),
		)),
	)
	return tp
}
