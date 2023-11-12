package handlers

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(handlers *Handlers) *chi.Mux {

	r := chi.NewRouter()
	m := NewMiddleware(handlers.logger)

	r.Use(m.WithLogging, m.WithCompress)

	r.Get("/{shortkey}", handlers.mainHandlerGet)
	r.Post("/", handlers.mainHandlerPost)

	r.Post("/api/shorten", handlers.APIShortenHandlerPost)
	r.Post("/api/shorten/batch", handlers.APIShortenBatchHandlerPost)

	r.Get("/ping", handlers.PingGet)

	return r
}
