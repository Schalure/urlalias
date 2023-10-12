package handlers

import "github.com/go-chi/chi/v5"

type Router struct {
	handlers Handlers
}

func NewRouter(handlers *Handlers) *chi.Mux {

	r := chi.NewRouter()

	r.Get("/{shortkey}", handlers.mainHandlerGet)
	r.Post("/", handlers.mainHandlerPost)

	return r
}
