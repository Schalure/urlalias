package aliasmaker

import (
	"fmt"
	"strings"

	"github.com/Schalure/urlalias/internal/app/storage"
)

const (
	aliasKeyLen        int = 9
)

// Type of service
type AliasMakerServise struct {
	Storage Storager
	lastKey string
}

// --------------------------------------------------
//	Constructor

func NewAliasMakerServise(storage Storager) *AliasMakerServise {
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

	newAliasKey, err := s.createAliasKey()
	if err != nil{
		return nil, err
	}

	return &storage.AliasURLModel{
		LongURL:  longURL,
		ShortKey: newAliasKey,
	}, nil
}


// --------------------------------------------------
//
//	Make short alias from URL
//	Output:
//		alias string - short alias to "longURL"
func (s *AliasMakerServise) createAliasKey() (string, error) {

	var charset []string = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", 
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", 
	}

	if s.lastKey == ""{
		s.lastKey = "000000000"
		return s.lastKey, nil
	}

	newKey := strings.Split(s.lastKey, "")
	if(len(newKey) != aliasKeyLen){
		return "", fmt.Errorf("a non-valid key was received from the repository: %s", s.lastKey)
	}

	for i := 0; i < aliasKeyLen; i++ {
		for n, char := range charset {
			if (newKey[i] == char) {
				if n == len(charset) - 1 {
					newKey[i] = charset[0]
					break;					
				} else {
					newKey[i] = charset[n + 1]
					s.lastKey = strings.Join(newKey, "")
					return s.lastKey, nil
				}
			}
		}
	}
	return "", fmt.Errorf("it is impossible to generate a new string because the storage is full")
}