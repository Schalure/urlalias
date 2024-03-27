// server package implementation methods of Server type
package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server struct for server object
type Server struct {
	router http.Handler
	server *http.Server
}

// Constructor of Handler type
func New(host string, handler *Handler, midleware *Middleware) *Server {

	r := chi.NewRouter()

	r.Use(midleware.WithLogging, midleware.WithCompress)

	r.Get("/{shortkey}", handler.redirect)
	r.Get("/ping", handler.PingGet)

	r.Group(func(r chi.Router) {

		r.Use(midleware.WithAuthentication)
		r.Post("/", handler.getShortURL)
		r.Post("/api/shorten", handler.apiGetShortURL)
		r.Post("/api/shorten/batch", handler.apiGetBatchShortURL)
	})

	r.Group(func(r chi.Router) {

		r.Use(midleware.WithVerification)
		r.Get("/api/user/urls", handler.apiGetUserAliases)
		r.Delete("/api/user/urls", handler.aipDeleteUserAliases)
	})

	return &Server{
		router: r,
		server: &http.Server{Addr: host, Handler: r},
	}
}

// Run starts server in HTTP/HTTPS mode
func (s *Server) Run(isHTTPS bool) error {

	if isHTTPS {
		return s.server.ListenAndServeTLS("", "")
	} else {
		return s.server.ListenAndServe()
	}
}

// Stop stops server
func (s *Server) Stop(ctx context.Context) error {
	log.Println("server shutdown...")
	return s.server.Shutdown(ctx)
}
