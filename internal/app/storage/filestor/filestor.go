package filestor

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/Schalure/urlalias/internal/app/storage"
)

type FileStorage struct {
	fileName string
	lastKey  string
	lastID   uint64
}

// ------------------------------------------------------------
//
//	FileStorage constructor
//	Output:
//		*FileStorage
func NewFileStorage(fileName string) (*FileStorage, error) {

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lastKey string
	var lastID uint64

	for i := 0; scanner.Scan(); i++ {
		var node storage.AliasURLModel
		if err := json.Unmarshal([]byte(scanner.Text()), &node); err != nil {
			return nil, errors.New("invalid file format")
		}

		lastID = node.ID
		lastKey = node.ShortKey
	}

	return &FileStorage{
		fileName: fileName,
		lastKey:  lastKey,
		lastID:   lastID,
	}, nil
}

// ------------------------------------------------------------
//
//	Save pair "shortKey, longURL" to db
//	This is interfase method of "Storager" interface
//	Input:
//		urlAliasNode *repositories.AliasURLModel
//	Output:
//		error - if not nil, can not save "urlAliasNode" because duplicate key
func (s *FileStorage) Save(urlAliasNode *storage.AliasURLModel) error {

	var data []byte
	file, err := os.OpenFile(s.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
func (s *FileStorage) SaveAll(urlAliasNodes []storage.AliasURLModel) error {

	var data []byte
	file, err := os.OpenFile(s.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
func (s *FileStorage) FindByShortKey(shortKey string) *storage.AliasURLModel {

	file, err := os.OpenFile(s.fileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		var node storage.AliasURLModel
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
func (s *FileStorage) FindByLongURL(longURL string) *storage.AliasURLModel {

	file, err := os.OpenFile(s.fileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		var node storage.AliasURLModel
		if err := json.Unmarshal([]byte(scanner.Text()), &node); err != nil {
			return nil
		}

		if longURL == node.LongURL {
			return &node
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
func (s *FileStorage) GetLastShortKey() string {
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
func (s *FileStorage) IsConnected() bool {
	return true
}

// ------------------------------------------------------------
//
//	Close connection to DB
//	This is interfase method of "Storager" interface
//	Output:
//		error
func (s *FileStorage) Close() error {
	return nil
}
