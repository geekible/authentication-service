package config

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
)

func initServiceMux() *chi.Mux {
	requestLogger := httplog.NewLogger("authentication-service", httplog.Options{
		JSON:     true,
		Concise:  true,
		LogLevel: "debug",
	})

	mux := chi.NewMux()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	mux.Use(httplog.RequestLogger(requestLogger))
	mux.Use(middleware.Compress(5, "application/json"))
	mux.Use(middleware.AllowContentType("application/json", "text/xml"))
	mux.Use(middleware.NoCache)
	mux.Use(middleware.StripSlashes)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	return mux
}
