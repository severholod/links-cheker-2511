package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"links-cheker-2511/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTPServer struct {
	HttpHandlers *HTTPHandlers
	Config       *config.Config
}

func NewHTTPServer(handlers *HTTPHandlers, config *config.Config) *HTTPServer {
	return &HTTPServer{
		HttpHandlers: handlers,
		Config:       config,
	}
}

func (server *HTTPServer) Start(log *slog.Logger) error {
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // работа с middleware и router
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/links", server.HttpHandlers.HandleSaveUrls)
	router.Post("/report", server.HttpHandlers.HandleGetUrls)

	srv := &http.Server{
		Addr:         server.Config.Address,
		Handler:      router,
		ReadTimeout:  server.Config.HTTPServer.Timeout,
		WriteTimeout: server.Config.HTTPServer.Timeout,
		IdleTimeout:  server.Config.HTTPServer.IdleTimeout,
	}

	serverErr := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		log.Info("Shutting down server gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		log.Info("Server stopped successfully")
		return nil

	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}
}
