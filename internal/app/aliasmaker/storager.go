package aliasmaker

import (
	"github.com/Schalure/urlalias/internal/app/models"
)

// Access interface to storage
type Storager interface {
	Save(urlAliasNode *models.AliasURLModel) error
	SaveAll(urlAliasNode []models.AliasURLModel) error
	FindByShortKey(shortKey string) *models.AliasURLModel
	FindByLongURL(longURL string) *models.AliasURLModel
	GetLastShortKey() string
	IsConnected() bool
	Close() error
}
