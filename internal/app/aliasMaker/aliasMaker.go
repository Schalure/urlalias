package aliasmaker

import (
	"math/rand"
	"net/url"

	"github.com/Schalure/urlalias/internal/app/config"
	"github.com/Schalure/urlalias/internal/app/database"
)

const(
	alias_len int = 9
	trys_to_make_alias int = 5
)

//--------------------------------------------------
//	Get short alias from URL
//	Input:
//		longURL string - URL
//	Output:
//		alias string - short alias to "longURL"
//		err error -
func GetAliasUrl(longURL string) (string, error){

	//	Check to valid URL
	if _, err := url.ParseRequestURI(longURL); err != nil{
		return "", err
	}


	aliasKey, ok := database.GetAliasKey(longURL)
	if !ok{
		var err error = nil
		//	try to make URL
		for i := 0; i < trys_to_make_alias + 1; i++{

			aliasKey = createAliasKey(longURL)

			if err := database.SavePair(longURL, aliasKey); err == nil{
				break;
			}
			if i == trys_to_make_alias{
				return "", err
			}
		}
	}
	return "http://" + config.HOST + aliasKey, nil
}

//--------------------------------------------------
//	Make short alias from URL
//	Input:
//		longURL string - URL
//	Output:
//		alias string - short alias to "longURL"
func createAliasKey(longURL string) string{
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	alliasKey := make([]byte, alias_len)

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
