// Application for URL shortening
package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/handlers"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// ------------------------------------------------------------
//
//	Main function
func main() {

	conf := config.NewConfig()

	aliasLogger, err := zap.NewDevelopment()
    if err != nil {
        // вызываем панику, если ошибка
        panic("cannot initialize zap")
    }
	defer aliasLogger.Sync()
	suggarLogger := aliasLogger.Sugar()


	storage := memstor.NewMemStorage()

	router := handlers.NewRouter(handlers.NewHandlers(storage, conf, suggarLogger))

	suggarLogger.Infow(fmt.Sprintf(
		"%s service have been started...", config.AppName),
		"Server address", conf.Host(),
		"Base URL", conf.BaseURL(),
		"Save log to file", conf.LogToFile(),
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
func NewLogger(c *config.Configuration) *slog.Logger {

	var l *slog.Logger

	if c.LogToFile() {
		panic("Logging to file no inplemented!!!")
	} else {
		l = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return l
}
