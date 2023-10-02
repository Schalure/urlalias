// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/internal/app/config"
	"github.com/Schalure/urlalias/internal/app/handlers"
)

// ------------------------------------------------------------
//	Main function
func main() {

	fmt.Printf("%s service have been started...\n", config.AppName)

	mux := RegistreHandlers(handlers.HandlersList)

	//	Run server
	log.Fatal(run(mux))
}

// ------------------------------------------------------------
//	Servise run.	
//	Input:
//		mux *http.ServeMux
//	Output:
//		err error - if servise have become panic or fatal error
func run(mux *http.ServeMux) error{
	return http.ListenAndServe(config.Host, mux)
}

// ------------------------------------------------------------
//	Add handlers to mux	
//	Input:
//		handlersList map[string] http.HandlerFunc - list of handlers signatur and handler functions
//	Output:
//		mux *http.ServeMux - handler router
func RegistreHandlers(handlersList map[string] http.HandlerFunc) *http.ServeMux{

	mux := http.NewServeMux()

	for k, v := range handlersList{
		mux.HandleFunc(k,v)
	}
	return mux
}