package aliasmaker

import (
	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/models"
	"github.com/Schalure/urlalias/internal/app/storage/filestor"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/Schalure/urlalias/internal/app/storage/postgrestor"
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

// ------------------------------------------------------------
//
//	New storage
//	Input:
//		storageType string
//	Output:
//		Storager
func NewStorage(c *config.Configuration) (Storager, error) {

	switch c.StorageType() {
	case config.DataBaseStor:
		return postgrestor.NewPostgreStor(c.DBConnection())
	case config.FileStor:
		return filestor.NewFileStorage(c.StorageFile())
	default:
		return memstor.NewMemStorage()
	}
}
