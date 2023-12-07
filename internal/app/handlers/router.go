package handlers

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(handlers *Handlers) *chi.Mux {

	r := chi.NewRouter()
	m := NewMiddleware(handlers.service)

	r.Use(m.WithLogging, m.WithCompress)

	r.Get("/{shortkey}", handlers.mainHandlerGet)
	r.Get("/ping", handlers.PingGet)

	r.Group(func(r chi.Router) {

		r.Use(m.WithAuthentication)
		r.Post("/", handlers.mainHandlerPost)
		r.Post("/api/shorten", handlers.APIShortenHandlerPost)
		r.Post("/api/shorten/batch", handlers.APIShortenBatchHandlerPost)
	})

	r.Group(func(r chi.Router) {

		r.Use(m.WithVerification)
		r.Get("/api/user/urls", handlers.APIUserURLsHandlerGet)
		r.Delete("/api/user/urls", handlers.APIUserURLsHandlerDelete)
	})

	return r
}
