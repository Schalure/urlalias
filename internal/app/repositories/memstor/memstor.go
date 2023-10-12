/*
This package describes types and methods for storing long URLs
and their alias keys in program memory.

Type "MemStorage" implements the "RepositoryURL" interface.
*/
package memstor

import (
	"fmt"

	"github.com/Schalure/urlalias/internal/app/repositories"
)

// Type for storage long URL and their alias keys
type MemStorage struct {
	stor map[string]string
}

// ------------------------------------------------------------
//
//	MemStorage constructor
//	Output:
//		*MemStorage
func NewMemStorage() *MemStorage {

	var s MemStorage
	s.stor = make(map[string]string)
	return &s
}

// ------------------------------------------------------------
//
//	Save pair "shortKey, longURL" to db
//	This is interfase method of "RepositoryURL" interface
//	Input:
//		urlAliasNode *repositories.AliasURLModel
//	Output:
//		error - if not nil, can not save "urlAliasNode" because duplicate key
func (s *MemStorage) Save(urlAliasNode *repositories.AliasURLModel) error {

	if _, ok := s.stor[urlAliasNode.ShortKey]; ok {
		return fmt.Errorf("the key \"%s\" is already in the database", urlAliasNode.ShortKey)
	}

	s.stor[urlAliasNode.ShortKey] = urlAliasNode.LongURL
	return nil
}

// ------------------------------------------------------------
//
//	Find "urlAliasNode models.AliasURLModel" by short key
//	This is interfase method of "RepositoryURL" interface
//	Input:
//		shortKey string
//	Output:
//		*repositories.AliasURLModel
//		error - if can not find "urlAliasNode" by short key
func (s *MemStorage) FindByShortKey(shortKey string) (*repositories.AliasURLModel, error) {

	longURL, ok := s.stor[shortKey]
	if !ok {
		return nil, fmt.Errorf("the urlAliasNode not found by key \"%s\"", shortKey)
	}
	return &repositories.AliasURLModel{ID: 0, ShortKey: shortKey, LongURL: longURL}, nil
}

// ------------------------------------------------------------
//
//	Find "urlAliasNode models.AliasURLModel" by long URL
//	This is interfase method of "RepositoryURL" interface
//	Input:
//		longURL string
//	Output:
//		*repositories.AliasURLModel
//		error - if can not find "urlAliasNode" by long URL
func (s *MemStorage) FindByLongURL(longURL string) (*repositories.AliasURLModel, error) {

	for k, v := range s.stor {
		if v == longURL {
			return &repositories.AliasURLModel{ID: 0, ShortKey: k, LongURL: longURL}, nil
		}
	}
	return nil, fmt.Errorf("the urlAliasNode not found by long URL \"%s\"", longURL)
}
