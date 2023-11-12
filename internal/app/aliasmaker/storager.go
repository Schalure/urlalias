package aliasmaker

import "github.com/Schalure/urlalias/internal/app/storage"

// Access interface to storage
type Storager interface {
	Save(urlAliasNode *storage.AliasURLModel) error
	SaveAll(urlAliasNode []storage.AliasURLModel) error
	FindByShortKey(shortKey string) *storage.AliasURLModel
	FindByLongURL(longURL string) *storage.AliasURLModel
	GetLastShortKey() string
	IsConnected() bool
	Close() error
}
