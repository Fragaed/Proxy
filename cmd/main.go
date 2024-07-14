package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"os/signal"
	"proxy/config"
	"proxy/internal/infrastructure/db/postgres"
	"proxy/run"
	"syscall"
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
