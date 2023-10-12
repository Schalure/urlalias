package handlers

import (
	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/repositories"
)

type Handlers struct {
	storage *repositories.RepositoryURL
	config  *config.Configuration
}

func NewHandlers(storage repositories.RepositoryURL, config *config.Configuration) *Handlers {

	return &Handlers{
		storage: &storage,
		config:  config,
	}
}
