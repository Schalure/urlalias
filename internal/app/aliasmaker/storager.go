package aliasmaker

import "github.com/Schalure/urlalias/internal/app/storage"

// Access interface to storage
type Storager interface {
	Save(urlAliasNode *storage.AliasURLModel) error
	FindByShortKey(shortKey string) (*storage.AliasURLModel, error)
	FindByLongURL(longURL string) (*storage.AliasURLModel, error)
}
