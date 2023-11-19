// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/handlers"
)

// ------------------------------------------------------------
//
//	Main function
func main() {

	conf := config.NewConfig()

	service, err := aliasmaker.NewAliasMakerServise(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Stop()

	router := handlers.NewRouter(handlers.NewHandlers(service))

	service.Logger.Infow(
		fmt.Sprintf("%s service have been started...", config.AppName),
		"Server address", conf.Host(),
		"Base URL", conf.BaseURL(),
		"Save log to file", conf.LogToFile(),
		"Storage file", conf.AliasesFile(),
		"DB connection string", conf.DBConnection(),
		"Storage type", conf.StorageType().String(),
	)

	err = http.ListenAndServe(conf.Host(), router)
	service.Logger.Fatalw(
		"aliasURL service stoped!",
		"error", err,
	)
}
