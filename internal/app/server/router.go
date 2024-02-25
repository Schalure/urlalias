package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *Server) http.Handler /*chi.Mux*/ {

	r := chi.NewRouter()
	m := NewMiddleware(handler.userManager, handler.logger)

	r.Use(m.WithLogging, m.WithCompress)

	r.Get("/{shortkey}", handler.redirect)
	r.Get("/ping", handler.PingGet)

	r.Group(func(r chi.Router) {

		r.Use(m.WithAuthentication)
		r.Post("/", handler.getShortURL)
		r.Post("/api/shorten", handler.apiGetShortURL)
		r.Post("/api/shorten/batch", handler.apiGetBatchShortURL)
	})

	r.Group(func(r chi.Router) {

		r.Use(m.WithVerification)
		r.Get("/api/user/urls", handler.apiGetUserAliases)
		r.Delete("/api/user/urls", handler.aipDeleteUserAliases)
	})

	return r
}
