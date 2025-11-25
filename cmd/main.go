package main

import (
	"links-cheker-2511/internal/config"
	"links-cheker-2511/internal/http"
	"links-cheker-2511/internal/storage"
	"log/slog"
	"os"
)

const (
	envProd = "prod"
	envDev  = "dev"
)

func main() {
	cfg := config.MustLoad()    // работа с конфигурациями
	log := setupLogger(cfg.Env) // работа с логами
	log.Info("starting links-checker", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	sqliteStorage, err := storage.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("error opening sqlite storage", err)
		os.Exit(1)
	}

	handlers := http.NewHTTPHandlers(sqliteStorage)
	httpServer := http.NewHTTPServer(handlers, cfg)

	log.Info("starting server", slog.String("address", cfg.Address))
	if err := httpServer.Start(); err != nil {
		log.Error("failed to start server", err)
	}

	log.Error("stopped server")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
