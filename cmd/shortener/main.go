// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/handlers"
	"github.com/Schalure/urlalias/internal/app/storage/filestor"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/Schalure/urlalias/internal/app/storage/postgrestor"
	"github.com/go-chi/chi/v5"
)

// ------------------------------------------------------------
//
//	Main function
func main() {

	conf := config.NewConfig()

	aliasLogger, err := handlers.NewLogger(handlers.LoggerTypeZap)
	if err != nil {
		log.Panicf("cannot initialize logger: %s", err)
	}
	defer aliasLogger.Close()

	stor := NewStorage(conf)
	defer stor.Close()

	service := aliasmaker.NewAliasMakerServise(stor)

	router := handlers.NewRouter(handlers.NewHandlers(service, conf, aliasLogger))

	aliasLogger.Infow(fmt.Sprintf(
		"%s service have been started...", config.AppName),
		"Server address", conf.Host(),
		"Base URL", conf.BaseURL(),
		"Save log to file", conf.LogToFile(),
		"Storage file", conf.StorageFile(),
		"DB connection string", conf.DBConnection(),
		"Storage type", conf.StorageType().String(),
	)

	log.Fatal(run(conf.Host(), router))
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

// ------------------------------------------------------------
//
//	New storage
//	Input:
//		storageType string
//	Output:
//		Storager
func NewStorage(c *config.Configuration) aliasmaker.Storager {

	switch c.StorageType() {
	case config.DataBaseStor:
		return postgrestor.NewPostgreStor(c.DBConnection())
	case config.FileStor:
		return filestor.NewFileStorage(c.StorageFile())
	default:
		return memstor.NewMemStorage()
	}
}
