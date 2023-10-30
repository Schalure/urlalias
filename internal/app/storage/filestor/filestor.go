package filestor

import "github.com/Schalure/urlalias/internal/app/storage"

type FileStorage struct {}

 
func (s *FileStorage) Save(urlAliasNode *storage.AliasURLModel) error{

}

func (s *FileStorage) FindByShortKey(shortKey string) (*storage.AliasURLModel, error){

}

func (s *FileStorage) FindByLongURL(longURL string) (*storage.AliasURLModel, error){
	
}
