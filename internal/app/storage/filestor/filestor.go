package filestor

import "github.com/Schalure/urlalias/internal/app/storage"

type FileStorage struct{
	stor []storage.AliasURLModel
}

// ------------------------------------------------------------
//
//	FileStorage constructor
//	Output:
//		*FileStorage
func NewFileStorage() *FileStorage{

	return &FileStorage{
		stor: make([]storage.AliasURLModel, 0),
	}
}

// ------------------------------------------------------------
//
//	Save pair "shortKey, longURL" to db
//	This is interfase method of "RepositoryURL" interface
//	Input:
//		urlAliasNode *repositories.AliasURLModel
//	Output:
//		error - if not nil, can not save "urlAliasNode" because duplicate key
func (s *FileStorage) Save(urlAliasNode *storage.AliasURLModel) error{

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
func (s *FileStorage) FindByShortKey(shortKey string) (*storage.AliasURLModel, error){
	return nil, nil

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
func (s *FileStorage) FindByLongURL(longURL string) (*storage.AliasURLModel, error){
	return nil, nil

}