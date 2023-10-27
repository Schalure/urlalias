package aliasmaker

import (
	"fmt"
	"math/rand"

	"github.com/Schalure/urlalias/internal/app/storage"
)

const (
	aliasKeyLen        int = 9
	trysToMakeAliasKey int = 5
)

//	Type of service
type AliasMakerServise struct{
	Storage Storager
}

// --------------------------------------------------
//	Constructor

func NewAliasMakerServise(storage Storager) *AliasMakerServise{
	return &AliasMakerServise{
		Storage: storage,
	}
}

// --------------------------------------------------
//
//	Create new URL pair
//	Output:
//		alias string - short alias to "longURL"
func (s *AliasMakerServise) NewPairURL(longURL string) (*storage.AliasURLModel, error) {

	for i := 0; i < trysToMakeAliasKey+1; i++ {

		node := storage.AliasURLModel{
			ID:       0,
			LongURL:  longURL,
			ShortKey: createAliasKey(),
		}

		if err := s.Storage.Save(&node); err == nil {
			return &node, nil
		}
	}
	return nil, fmt.Errorf("can not create alias key from \"%s\"", longURL)
}

// --------------------------------------------------
//
//	Make short alias from URL
//	Output:
//		alias string - short alias to "longURL"
func createAliasKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	alliasKey := make([]byte, aliasKeyLen)

	for i := range alliasKey {
		alliasKey[i] = charset[rand.Intn(len(charset))]
	}

	return string(alliasKey)
}
