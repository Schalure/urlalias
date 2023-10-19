// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/handlers"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/go-chi/chi/v5"
)

// ------------------------------------------------------------
//
//	Main function
func main() {

	fmt.Printf("%s service have been started...\n", config.AppName)

	config := config.NewConfig()

	storage := memstor.NewMemStorage()

	router := handlers.NewRouter(handlers.NewHandlers(storage, config))

	log.Fatal(run(config.Host(), router))
}

// ------------------------------------------------------------
//
//	Servise run.
//	Input:
//		mux *chi.Mux
//	Output:
//		err error - if servise have become panic or fatal error
func run(serverAddres string, router *chi.Mux) error {
	return http.ListenAndServe(serverAddres, router)
}
