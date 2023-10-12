package handlers

import (
	"github.com/Schalure/urlalias/models"
	"github.com/go-chi/chi/v5"
)

type Handlers struct{
	storege models.RepositoryURL
}

func NewHandlers(storage models.RepositoryURL) *Handlers{
	return &Handlers{
		storege: storage,
	}
}

type Router struct{
	handlers Handlers
}

func NewRouter(handlers Handlers) *chi.Mux{

	r := chi.NewRouter()

	r.Get("/{shortkey}", handlers.mainHandlerGet)
	r.Post("/", handlers.mainHandlerPost)

	return r
}