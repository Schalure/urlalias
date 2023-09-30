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
//	Make short alias from URL
//	Input:
//		longURL string - URL
//	Output:
//		alias string - short alias to "longURL"
//		err error -
func MakeAliasUrl(longURL string) (string, error){

	//	Check to valid URL
	if _, err := url.ParseRequestURI(longURL); err != nil{
		return "", err
	}

	alias, ok := database.GetAliasFromDB(longURL)
	if !ok{
		var err error = nil
		//	try to make URL
		for i := 0; i < trys_to_make_alias; i++{
			alias = createAlias(longURL)
			if err = database.SavePairToDB(longURL, alias); err == nil{
				break;
			}
		}
		//	if can't create alias
		if err != nil{
			return "", err
		}
	}
	return alias, nil
}

//--------------------------------------------------
//	Make short alias from URL
//	Input:
//		longURL string - URL
//	Output:
//		alias string - short alias to "longURL"
func createAlias(longURL string) string{
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	allias := make([]byte, alias_len)

	for i := range allias{
		allias[i] = charset[rand.Intn(len(charset))]
	}

	return "http://" + config.HOST + "/" + string(allias)
}
