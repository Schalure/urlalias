package filestor

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Schalure/urlalias/internal/app/storage"
)

type FileStorage struct {
	stor     map[string]string
	fileName string
}

// ------------------------------------------------------------
//
//	FileStorage constructor
//	Output:
//		*FileStorage
func NewFileStorage(fileName string) (*FileStorage, error) {

	return &FileStorage{
		stor:     make(map[string]string),
		fileName: fileName,
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

	if _, ok := s.stor[urlAliasNode.ShortKey]; ok {
		return fmt.Errorf("the key \"%s\" is already in the database", urlAliasNode.ShortKey)
	}

	s.stor[urlAliasNode.ShortKey] = urlAliasNode.LongURL
	urlAliasNode.ID = uint64(len(s.stor))

	var data []byte
	file, err := os.OpenFile(s.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if data, err = json.Marshal(urlAliasNode); err != nil {
		return err
	}

	if _, err = file.Write(append(data, '\n')); err != nil {
		return err
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

	longURL, ok := s.stor[shortKey]
	if !ok {
		return nil
	}
	return &storage.AliasURLModel{ID: 0, ShortKey: shortKey, LongURL: longURL}
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

	for k, v := range s.stor {
		if v == longURL {
			return &storage.AliasURLModel{ID: 0, ShortKey: k, LongURL: longURL}
		}
	}
	return nil
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
