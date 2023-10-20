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
)

// ------------------------------------------------------------
//
//	Main function
func main() {

	c := config.NewConfig()

	aliasLogger := NewLogger(c)

	storage := memstor.NewMemStorage()

	router := handlers.NewRouter(handlers.NewHandlers(storage, c, aliasLogger))

	aliasLogger.Info(fmt.Sprintf(
		"%s service have been started...", config.AppName),
		"Server address", c.Host(),
		"Base URL", c.BaseURL(),
		"Save log to file", c.LogToFile(),
	)

	log.Fatal(run(c.Host(), router))
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
