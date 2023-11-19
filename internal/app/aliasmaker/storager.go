package aliasmaker

import (
	"context"

	"github.com/Schalure/urlalias/internal/app/models"
)

// Access interface to storage
type Storager interface {
	CreateUser() (uint64, error)
	Save(urlAliasNode *models.AliasURLModel) error
	SaveAll(urlAliasNode []models.AliasURLModel) error
	FindByShortKey(shortKey string) *models.AliasURLModel
	FindByLongURL(longURL string) *models.AliasURLModel
	FindByUserID(ctx context.Context, userID uint64) ([]models.AliasURLModel, error)
	GetLastShortKey() string
	IsConnected() bool
	Close() error
}
