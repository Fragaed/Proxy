package main

import (
	"flag"
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
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Info("Failed to sync logger", zap.Error(err))
		}
	}()

	// Загружаем конфигурацию по умолчанию
	cfg := config.MustLoad()
	logger.Infow("Конфигурация загружена по умолчанию")

	// Определение флагов для переопределения конфигурации базы данных
	dbHost := flag.String("db-host", cfg.DB.Host, "Хост базы данных")
	dbPort := flag.String("db-port", cfg.DB.Port, "Порт базы данных")
	dbUser := flag.String("db-user", cfg.DB.Username, "Пользователь базы данных")
	dbPassword := flag.String("db-password", cfg.DB.Password, "Пароль базы данных")
	dbName := flag.String("db-name", cfg.DB.DBName, "Имя базы данных")

	// Разбор флагов
	flag.Parse()

	// Переопределение конфигурации базы данных, если флаги были заданы
	cfg.DB.Host = *dbHost
	cfg.DB.Port = *dbPort
	cfg.DB.Username = *dbUser
	cfg.DB.Password = *dbPassword
	cfg.DB.DBName = *dbName

	logger.Infow("Конфигурация базы данных загружена", "host", cfg.DB.Host, "port", cfg.DB.Port, "user", cfg.DB.Username, "db", cfg.DB.DBName)

	// Инициализация подключения к базе данных
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
