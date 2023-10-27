package handlers

import (
	"github.com/Schalure/urlalias/internal/app/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handlers *Handlers) *chi.Mux {

	r := chi.NewRouter()
	m := middleware.NewMiddleware(handlers.logger)

	r.Use(m.WhithLogging)

	r.Get("/{shortkey}", handlers.mainHandlerGet)
	r.Post("/", handlers.mainHandlerPost)

	r.Post("/api/shorten", handlers.APIShortenHandlerPost)

	return r
}
