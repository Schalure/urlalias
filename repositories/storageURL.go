package repositories

import (
	"fmt"

	"github.com/Schalure/urlalias/models"
)

type StorageURL struct {
	db map[string]string
}

func NewStorageURL() *StorageURL {

	var s StorageURL
	s.db = make(map[string]string)
	return &s
}

// ------------------------------------------------------------
//
//	Save pair "shortKey, longURL" to db
//	This is interfase method of "RepositoryURL" interface
//	Input:
//		urlAliasNode models.AliasURLModel
//	Output:
//		error - if not nil, can not save "urlAliasNode" because duplicate key
func (s *StorageURL) Save(urlAliasNode models.AliasURLModel) (*models.AliasURLModel, error) {

	if _, ok := s.db[urlAliasNode.ShortKey]; ok {
		return nil, fmt.Errorf("the key \"%s\" is already in the database", urlAliasNode.ShortKey)
	}

	s.db[urlAliasNode.ShortKey] = urlAliasNode.LongURL
	return s.FindByShortKey(urlAliasNode.ShortKey)
}

// ------------------------------------------------------------
//
//	Find "urlAliasNode models.AliasURLModel" by short key
//	This is interfase method of "RepositoryURL" interface
//	Input:
//		shortKey string
//	Output:
//		*models.AliasURLModel
//		error - if can not find "urlAliasNode" by short key
func (s *StorageURL) FindByShortKey(shortKey string) (*models.AliasURLModel, error) {

	longURL, ok := s.db[shortKey]
	if !ok {
		return nil, fmt.Errorf("the urlAliasNode not found by key \"%s\"", shortKey)
	}
	return &models.AliasURLModel{ID: 0, ShortKey: shortKey, LongURL: longURL}, nil
}

// ------------------------------------------------------------
//
//	Find "urlAliasNode models.AliasURLModel" by long URL
//	This is interfase method of "RepositoryURL" interface
//	Input:
//		longURL string
//	Output:
//		*models.AliasURLModel
//		error - if can not find "urlAliasNode" by long URL
func (s *StorageURL) FindByLongURL(longURL string) (*models.AliasURLModel, error) {

	for k, v := range s.db {
		if v == longURL {
			return &models.AliasURLModel{ID: 0, ShortKey: k, LongURL: longURL}, nil
		}
	}
	return nil, fmt.Errorf("the urlAliasNode not found by long URL \"%s\"", longURL)
}
