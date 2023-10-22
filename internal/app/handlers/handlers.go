package handlers

import (
	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/storage"
	"go.uber.org/zap"
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
	logger  *zap.SugaredLogger
}

func NewHandlers(storage IStorage, config *config.Configuration, logger *zap.SugaredLogger) *Handlers {

	return &Handlers{
		storage: storage,
		config:  config,
		logger:  logger,
	}
}
