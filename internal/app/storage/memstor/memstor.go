/*
This package describes types and methods for storing long URLs
and their alias keys in program memory.

Type "MemStorage" implements the "RepositoryURL" interface.
*/
package memstor

import (
	"context"

	"github.com/Schalure/urlalias/internal/app/models"
)

// Type for storage long URL and their alias keys
type Storage struct {
	//	[key, value] = [ShortKey, LongURL]
	aliases []models.AliasURLModel
	users   []models.UserModel

	lastKey string
}

// ------------------------------------------------------------
//
//	MemStorage constructor
func NewStorage() (*Storage, error) {

	var s Storage
	s.aliases = make([]models.AliasURLModel, 0)
	s.users = make([]models.UserModel, 0)

	return &s, nil
}

// ------------------------------------------------------------
//
//	Create new user
func (s *Storage) CreateUser() (uint64, error) {

	user := models.UserModel{
		UserID: uint64(len(s.users)),
	}

	s.users = append(s.users, user)
	return user.UserID, nil
}

// ------------------------------------------------------------
//
//	Save pair "shortKey, longURL" to db
func (s *Storage) Save(urlAliasNode *models.AliasURLModel) error {

	s.aliases = append(s.aliases, *urlAliasNode)
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
func (s *Storage) SaveAll(urlAliasNodes []models.AliasURLModel) error {

	for _, node := range urlAliasNodes {

		s.aliases = append(s.aliases, node)
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
func (s *Storage) FindByShortKey(shortKey string) *models.AliasURLModel {

	for _, node := range s.aliases {
		if node.ShortKey == shortKey {
			return &node
		}
	}
	return nil
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
func (s *Storage) FindByLongURL(longURL string) *models.AliasURLModel {

	for _, node := range s.aliases {
		if node.LongURL == longURL {
			return &node
		}
	}
	return nil
}

// ------------------------------------------------------------
//
//	Find all "urlAliasNode models.AliasURLModel" by UserID
func (s *Storage) FindByUserID(ctx context.Context, userID uint64) ([]models.AliasURLModel, error) {

	var nodes []models.AliasURLModel

	for _, node := range s.aliases {
		if node.UserID == userID {
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}


// ------------------------------------------------------------
//
//	Mark aliases like "deleted" by aliasesID
func (s *Storage) MarkDeleted(ctx context.Context, aliasesID []uint64) error {

	return nil
}

// ------------------------------------------------------------
//
//	Get the last saved key
//	This is interfase method of "Storager" interface
//	Output:
//		string - last saved key
func (s *Storage) GetLastShortKey() string {
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
func (s *Storage) IsConnected() bool {
	return true
}

// ------------------------------------------------------------
//
//	Close connection to DB
//	This is interfase method of "Storager" interface
//	Output:
//		error
func (s *Storage) Close() error {
	return nil
}
