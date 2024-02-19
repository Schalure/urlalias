// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliaslogger"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/handlers"
	"github.com/Schalure/urlalias/internal/app/storage"
)

// ------------------------------------------------------------
//
//	Main function
func main() {

	log.Println("Start initialize application...")

	log.Println("Cofiguration initialize...")
	conf := config.NewConfig()

	log.Println("Logger initialize...")
	logger, err := aliaslogger.NewLogger(aliaslogger.LoggerTypeZap)
	if err != nil {
		log.Fatalln("Error, while initialization logger!", err)
	}

	log.Println("Storage initialize...")
	stor, err := storage.NewStorage(conf)
	if err != nil {
		log.Fatalln("Error, while initialization storage!", err)
	}

	log.Println("Alias maker service initialize...")
	service, err := aliasmaker.New(conf, stor, logger)
	if err != nil {
		log.Fatalln("Error, while initialization Alias maker service!", err)
	}
	defer service.Stop()

	log.Println("Router initialize...")
	router := handlers.NewRouter(handlers.New(service))

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
