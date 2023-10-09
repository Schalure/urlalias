// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/handlers"
	"github.com/Schalure/urlalias/repositories"
	"github.com/go-chi/chi/v5"
)

// ------------------------------------------------------------
//
//	Main function
func main() {

	fmt.Printf("%s service have been started...\n", config.AppName)

	//	Read application options
	config.Initialize()

	//	initialize storage
	storage := repositories.NewStorageURL()

	//	initialize router and handlers
	router := handlers.MakeRouter(storage)

	//	Run server
	log.Fatal(run(router))
}

// ------------------------------------------------------------
//
//	Servise run.
//	Input:
//		mux *chi.Mux
//	Output:
//		err error - if servise have become panic or fatal error
func run(router *chi.Mux) error {
	return http.ListenAndServe(config.Config.Host(), router)
}


