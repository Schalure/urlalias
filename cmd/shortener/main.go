// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/internal/app/config"
	"github.com/Schalure/urlalias/internal/app/handlers"
	"github.com/Schalure/urlalias/repositories"
	"github.com/go-chi/chi"
)

// ------------------------------------------------------------
//	Main function
func main() {

	fmt.Printf("%s service have been started...\n", config.AppName)

	router := RegistreHandlers()

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
	return http.ListenAndServe(config.Host, router)
}

// ------------------------------------------------------------
//	Add handlers to mux	
//	Input:
//		handlersList map[string] http.HandlerFunc - list of handlers signatur and handler functions
//	Output:
//		mux *http.ServeMux - handler router
func RegistreHandlers() *chi.Mux{

	//	create new router
	router := chi.NewRouter()

	//	create storage
	storage := repositories.NewStorageURL()

	router.Get("/", handlers.MainHandlerMethodGet(storage))
	router.Post("/{shortkey}", handlers.MainHandlerMethodPost(storage))

	return router
}