package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"links-cheker-2511/internal/config"
	"net/http"
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

func (server *HTTPServer) Start() error {
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // работа с middleware и router
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/links", func(w http.ResponseWriter, r *http.Request) {
		//server.Logger.Info("Hi from links checker", r.URL, r.Method, r.Body)
	})
	router.Post("/report", func(w http.ResponseWriter, r *http.Request) {
		//server.Logger.Info("Hi from links checker", r.URL, r.Method, r.Body)
	})

	srv := &http.Server{
		Addr:         server.Config.Address,
		Handler:      router,
		ReadTimeout:  server.Config.HTTPServer.Timeout,
		WriteTimeout: server.Config.HTTPServer.Timeout,
		IdleTimeout:  server.Config.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
