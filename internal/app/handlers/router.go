package handlers

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *Handler) *chi.Mux {

	r := chi.NewRouter()
	m := NewMiddleware(handler.service)

	r.Use(m.WithLogging, m.WithCompress)

	r.Get("/{shortkey}", handler.redirect)
	r.Get("/ping", handler.PingGet)

	r.Group(func(r chi.Router) {

		r.Use(m.WithAuthentication)
		r.Post("/", handler.getShortURL)
		r.Post("/api/shorten", handler.apiGetShortURL)
		r.Post("/api/shorten/batch", handler.APIShortenBatchHandlerPost)
	})

	r.Group(func(r chi.Router) {

		r.Use(m.WithVerification)
		r.Get("/api/user/urls", handler.APIUserURLsHandlerGet)
		r.Delete("/api/user/urls", handler.APIUserURLsHandlerDelete)
	})

	return r
}
