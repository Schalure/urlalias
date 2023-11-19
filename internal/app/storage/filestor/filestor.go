package filestor

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/Schalure/urlalias/internal/app/models"
)

type Storage struct {
	aliasesFileName string
	usersFileName   string
	lastKey         string
	lastID          uint64
	lastUserID      uint64
}

// ------------------------------------------------------------
//
//	FileStorage constructor
//	Output:
//		*FileStorage
func NewStorage(aliasesFileName, usersFileName string) (*Storage, error) {

	aliasesFile, err := os.OpenFile(aliasesFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer aliasesFile.Close()

	scanner := bufio.NewScanner(aliasesFile)

	var lastKey string
	var lastID uint64

	for i := 0; scanner.Scan(); i++ {
		var node models.AliasURLModel
		if err := json.Unmarshal([]byte(scanner.Text()), &node); err != nil {
			return nil, errors.New("invalid file format")
		}

		lastID = node.ID
		lastKey = node.ShortKey
	}

	usersFile, err := os.OpenFile(usersFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer usersFile.Close()

	scanner = bufio.NewScanner(usersFile)

	var lastUserID uint64

	for i := 0; scanner.Scan(); i++ {
		var node models.UserModel
		if err := json.Unmarshal([]byte(scanner.Text()), &node); err != nil {
			return nil, errors.New("invalid file format")
		}

		lastUserID = node.UserID
	}

	return &Storage{
		aliasesFileName: aliasesFileName,
		usersFileName:   usersFileName,
		lastKey:         lastKey,
		lastID:          lastID,
		lastUserID:      lastUserID,
	}, nil
}

// ------------------------------------------------------------
//
//	Create new user
func (s *Storage) CreateUser() (uint64, error) {

	var data []byte

	file, err := os.OpenFile(s.usersFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	newUserID := s.lastUserID + 1
	user := models.UserModel{
		UserID: newUserID,
	}

	if data, err = json.Marshal(user); err != nil {
		return 0, err
	}

	if _, err = file.Write(append(data, '\n')); err != nil {
		return 0, err
	}

	s.lastUserID = newUserID
	return s.lastUserID, nil
}

// ------------------------------------------------------------
//
//	Save pair "shortKey, longURL" to db
//	This is interfase method of "Storager" interface
//	Input:
//		urlAliasNode *repositories.AliasURLModel
//	Output:
//		error - if not nil, can not save "urlAliasNode" because duplicate key
func (s *Storage) Save(urlAliasNode *models.AliasURLModel) error {

	var data []byte
	file, err := os.OpenFile(s.aliasesFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	urlAliasNode.ID = s.lastID + 1
	if data, err = json.Marshal(urlAliasNode); err != nil {
		return err
	}

	if _, err = file.Write(append(data, '\n')); err != nil {
		return err
	}

	s.lastID++
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

	var data []byte
	file, err := os.OpenFile(s.aliasesFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, node := range urlAliasNodes {
		node.ID = s.lastID + 1
		if data, err = json.Marshal(&node); err != nil {
			return err
		}

		if _, err = file.Write(append(data, '\n')); err != nil {
			return err
		}

		s.lastID++
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

	file, err := os.OpenFile(s.aliasesFileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		var node models.AliasURLModel
		if err := json.Unmarshal([]byte(scanner.Text()), &node); err != nil {
			return nil
		}

		if shortKey == node.ShortKey {
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

	file, err := os.OpenFile(s.aliasesFileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		var node models.AliasURLModel
		if err := json.Unmarshal([]byte(scanner.Text()), &node); err != nil {
			return nil
		}

		if longURL == node.LongURL {
			return &node
		}

	}
	return nil
}

func (s *Storage) FindByUserID(ctx context.Context, userID uint64) ([]models.AliasURLModel, error) {

	file, err := os.OpenFile(s.aliasesFileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nodes []models.AliasURLModel
	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		var node models.AliasURLModel
		if err := json.Unmarshal([]byte(scanner.Text()), &node); err != nil {
			return nil, err
		}

		if node.UserID == userID {
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
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
