// Application for URL shortening
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/server"
	"github.com/Schalure/urlalias/internal/app/storage"
)

var (
	//	buildVersion is last build version of shortner service
	buildVersion string = "N/A"
	//	buildDate is last build date of shortner service
	buildDate string = "N/A"
	//	buildCommit is last commit of shortner service
	buildCommit string = "N/A"
)

// ------------------------------------------------------------
//
//	Main function
func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	log.Println("Start initialize application...")
	ctxStop, cancelStop := context.WithCancel(context.Background())
	defer cancelStop()

	log.Println("Cofiguration initialize...")
	conf := config.NewConfig()

	log.Println("Logger initialize...")
	logger, err := zaplogger.NewZapLogger("")
	if err != nil {
		log.Fatalln("Error, while initialization logger!", err)
	}

	log.Println("Storage initialize...")
	stor, err := storage.NewStorage(conf)
	if err != nil {
		log.Fatalln("Error, while initialization storage!", err)
	}

	log.Println("Alias maker service initialize...")
	service, err := aliasmaker.New(stor, logger)
	if err != nil {
		log.Fatalln("Error, while initialization Alias maker service!", err)
	}
	service.Run(ctxStop)
	defer service.Stop()

	log.Println("Router initialize...")
	router := server.NewRouter(server.New(service, service, logger, conf.BaseURL()))

	logger.Infow(
		fmt.Sprintf("%s service have been started...", config.AppName),
		"Server address", conf.Host(),
		"Base URL", conf.BaseURL(),
		"Save log to file", conf.LogToFile(),
		"Storage file", conf.AliasesFile(),
		"DB connection string", conf.DBConnection(),
		"Storage type", conf.StorageType().String(),
		"Is HTTPS", conf.EnableHTTPS(),
	)


	if conf.EnableHTTPS() {
		server := &http.Server{
			Addr:      ":443",
			Handler:   router,
		}
		err = server.ListenAndServeTLS("", "")
	} else {
		err = http.ListenAndServe(conf.Host(), router)
	}
	
	logger.Fatalw(
		"aliasURL service stoped!",
		"error", err,
	)
}
