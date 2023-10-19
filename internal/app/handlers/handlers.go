package handlers

import (
	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/storage"
)

// Access interface to storage
type IStorage interface {
	Save(s *storage.AliasURLModel) error
	FindByShortKey(shortKey string) (*storage.AliasURLModel, error)
	FindByLongURL(longURL string) (*storage.AliasURLModel, error)
}


type Handlers struct {
	storage IStorage
	config  *config.Configuration
}

func NewHandlers(storage IStorage, config *config.Configuration) *Handlers {

	return &Handlers{
		storage: storage,
		config:  config,
	}
}
