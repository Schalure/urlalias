/*
This package describes types and methods for storing long URLs
and their alias keys in program memory.

Type "MemStorage" implements the "RepositoryURL" interface.
*/
package memstor

import "github.com/Schalure/urlalias/internal/app/models"

// Type for storage long URL and their alias keys
type MemStorage struct {
	//	[key, value] = [ShortKey, LongURL]
	stor    map[string]string
	lastKey string
}

// ------------------------------------------------------------
//
//	MemStorage constructor
//	Output:
//		*MemStorage
func NewMemStorage() (*MemStorage, error) {

	var s MemStorage
	s.stor = make(map[string]string)

	return &s, nil
}

// ------------------------------------------------------------
//
//	Save pair "shortKey, longURL" to db
//	This is interfase method of "Storager" interface
//	Input:
//		urlAliasNode *repositories.AliasURLModel
//	Output:
//		error - if not nil, can not save "urlAliasNode" because duplicate key
func (s *MemStorage) Save(urlAliasNode *models.AliasURLModel) error {

	s.stor[urlAliasNode.ShortKey] = urlAliasNode.LongURL
	s.lastKey = urlAliasNode.ShortKey

	return nil
}

// ------------------------------------------------------------
//
//	Save array of pairs "shortKey, longURL" to db
//	This is interfase method of "Storager" interface
//	Input:
//		urlAliasNode []repositories.AliasURLModel
//	Output:
//		error - if not nil, can not save "[]storage.AliasURLModel"
func (s *MemStorage) SaveAll(urlAliasNodes []models.AliasURLModel) error {

	for _, node := range urlAliasNodes {
		s.stor[node.ShortKey] = node.LongURL
		s.lastKey = node.ShortKey
	}
	return nil
}

// ------------------------------------------------------------
//
//	Find "urlAliasNode models.AliasURLModel" by short key
//	This is interfase method of "Storager" interface
//	Input:
//		shortKey string
//	Output:
//		*repositories.AliasURLModel
//		error - if can not find "urlAliasNode" by short key
func (s *MemStorage) FindByShortKey(shortKey string) *models.AliasURLModel {

	longURL, ok := s.stor[shortKey]
	if !ok {
		return nil
	}
	return &models.AliasURLModel{ID: 0, ShortKey: shortKey, LongURL: longURL}
}

// ------------------------------------------------------------
//
//	Find "urlAliasNode models.AliasURLModel" by long URL
//	This is interfase method of "Storager" interface
//	Input:
//		longURL string
//	Output:
//		*repositories.AliasURLModel
//		error - if can not find "urlAliasNode" by long URL
func (s *MemStorage) FindByLongURL(longURL string) *models.AliasURLModel {

	for k, v := range s.stor {
		if v == longURL {
			return &models.AliasURLModel{ID: 0, ShortKey: k, LongURL: longURL}
		}
	}
	return nil
}

// ------------------------------------------------------------
//
//	Get the last saved key
//	This is interfase method of "Storager" interface
//	Output:
//		string - last saved key
func (s *MemStorage) GetLastShortKey() string {
	return s.lastKey
}

// ------------------------------------------------------------
//
//	Check connection to DB
//	This is interfase method of "Storager" interface
//	Output:
//		bool - true: connection is
//			   false: connection isn't
//		error - if can not find "urlAliasNode" by long URL
func (s *MemStorage) IsConnected() bool {
	return true
}

// ------------------------------------------------------------
//
//	Close connection to DB
//	This is interfase method of "Storager" interface
//	Output:
//		error
func (s *MemStorage) Close() error {
	return nil
}
