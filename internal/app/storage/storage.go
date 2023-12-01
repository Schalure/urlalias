package storage

import (
	"github.com/Schalure/urlalias/cmd/shortener/config"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/storage/filestor"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
	"github.com/Schalure/urlalias/internal/app/storage/postgrestor"
)

// --------------------------------------------------
//
//	Choose storage for service
func NewStorage(c *config.Configuration) (aliasmaker.Storager, error) {

	switch c.StorageType() {
	case config.DataBaseStor:
		return postgrestor.NewStorage(c.DBConnection())
	case config.FileStor:
		return filestor.NewStorage(c.AliasesFile(), c.UsersFile())
	default:
		return memstor.NewStorage()
	}
}