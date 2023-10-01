package aliasmaker

import (
	"math/rand"
	"net/url"

	"github.com/Schalure/urlalias/internal/app/config"
	"github.com/Schalure/urlalias/internal/app/storage"
)

const(
	aliasKeyLen int = 9
	trysToMakeAliasKey int = 5
)

//--------------------------------------------------
//	Get short alias from URL
//	Input:
//		longURL string - URL
//	Output:
//		alias string - short alias to "longURL"
//		err error -
func GetAliasURL(longURL string) (string, error){

	//	Check to valid URL
	if _, err := url.ParseRequestURI(longURL); err != nil{
		return "", err
	}


	aliasKey, ok := storage.GetAliasKey(longURL)
	if !ok{
		var err error = nil
		//	try to make URL
		for i := 0; i < trysToMakeAliasKey + 1; i++{

			aliasKey = createAliasKey(longURL)

			if err := storage.SavePair(longURL, aliasKey); err == nil{
				break;
			}
			if i == trysToMakeAliasKey{
				return "", err
			}
		}
	}
	return "http://" + config.Host + aliasKey, nil
}

//--------------------------------------------------
//	Make short alias from URL
//	Input:
//		longURL string - URL
//	Output:
//		alias string - short alias to "longURL"
func createAliasKey(longURL string) string{
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	alliasKey := make([]byte, aliasKeyLen)

	for i := range alliasKey{
		alliasKey[i] = charset[rand.Intn(len(charset))]
	}

	return "/" + string(alliasKey)
}

//--------------------------------------------------
//	Get short alias from URL
//	Input:
//		key string - short key
//	Output:
//		alias string - short alias to "longURL"
//func getAlias(key )
