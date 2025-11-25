package main

import (
	"links-cheker-2511/internal/config"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()
	router.Use(middleware.RequestID) // работа с middleware и router
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/links", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Hi from links checker", r.URL, r.Method, r.Body)
	})
	router.Post("/report", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Hi from links checker", r.URL, r.Method, r.Body)
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
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
