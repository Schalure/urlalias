// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/internal/app"
)

//--------------------------------------------------
//	Main function
func main() {

	fmt.Println("github.com/Schalure/urlalias service have been started...")

	mux := RegistreHandlers(app.HandlersList)

	//	Run server
	log.Fatal(run(mux))
}

//--------------------------------------------------
//	Servise run
func run(mux *http.ServeMux) error{
	return http.ListenAndServe("192.168.1.88:8080", mux)
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




