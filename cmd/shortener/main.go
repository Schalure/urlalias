// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/handlers"
	"github.com/Schalure/urlalias/repositories"
	"github.com/go-chi/chi"
)

// ------------------------------------------------------------
//	Main function
func main() {

	fmt.Printf("%s service have been started...\n", config.AppName)

	//	Read application options
	config.MustInit()

	//	initialize storage
	storage := repositories.NewStorageURL()

	//	initialize router and handlers
	router := RegistreHandlers(storage)

	//	Run server
	log.Fatal(run(router))
}

// ------------------------------------------------------------
//	Servise run.	
//	Input:
//		mux *http.ServeMux
//	Output:
//		err error - if servise have become panic or fatal error
func run(router chi.Router) error{
	return http.ListenAndServe(config.Config.Host, router)
}

// ------------------------------------------------------------
//	Add handlers to router	
//	Output:
//		router *chi.Mux - handler router
func RegistreHandlers(storage *repositories.StorageURL) *chi.Mux{

	//	create new router
	router := chi.NewRouter()

	//	add handlers
	router.Get("/{shortkey}", handlers.MainHandlerMethodGet(storage))
	router.Post("/", handlers.MainHandlerMethodPost(storage))

	return router
}