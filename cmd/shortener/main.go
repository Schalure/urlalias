// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/internal/app/config"
	"github.com/Schalure/urlalias/internal/app/handlers"
)

//--------------------------------------------------
//	Main function
func main() {

	fmt.Printf("%s service have been started...\n", config.AppName)

	mux := RegistreHandlers(handlers.HandlersList)

	//	Run server
	log.Fatal(run(mux))
}

//--------------------------------------------------
//	Servise run
func run(mux *http.ServeMux) error{
	return http.ListenAndServe(config.Host, mux)
}

//--------------------------------------------------
//	Add handlers to mux
func RegistreHandlers(handlersList map[string] http.HandlerFunc) *http.ServeMux{

	mux := http.NewServeMux()

	for k, v := range handlersList{
		mux.HandleFunc(k,v)
	}

	return mux
}




