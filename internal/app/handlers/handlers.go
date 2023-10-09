package handlers

import (
	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/repositories"
	"github.com/go-chi/chi/v5"
)

// ------------------------------------------------------------
//
//	Add handlers to router
//	Output:
//		router *chi.Mux - handler router
func MakeRouter(storage *repositories.StorageURL) *chi.Mux {

	//	create new router
	router := chi.NewRouter()

	//	add handlers
	router.Get("/{shortkey}", MainHandlerMethodGet(storage))
	router.Post("/", MainHandlerMethodPost(storage, config.Config))

	return router
}