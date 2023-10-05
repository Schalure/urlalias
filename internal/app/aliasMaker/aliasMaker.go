package aliasmaker

import (
	"math/rand"
)

const(
	aliasKeyLen int = 9
	TrysToMakeAliasKey int = 5
)

//--------------------------------------------------
//	Get short alias from URL
//	Input:
//		longURL string - URL
//	Output:
//		alias string - short alias to "longURL"
//		err error -
// func GetAliasURL(longURL string) (string, error){

// 	aliasKey, ok := storage.GetAliasKey(longURL)
// 	if !ok{
// 		var err error = nil
// 		//	try to make URL
// 		for i := 0; i < TrysToMakeAliasKey + 1; i++{

// 			aliasKey = createAliasKey(longURL)

// 			if err := storage.SavePair(longURL, aliasKey); err == nil{
// 				break;
// 			}
// 			if i == TrysToMakeAliasKey{
// 				return "", err
// 			}
// 		}
// 	}
// 	return "http://" + config.Host + aliasKey, nil
// }

//--------------------------------------------------
//	Make short alias from URL
//	Output:
//		alias string - short alias to "longURL"
func CreateAliasKey() string{
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	alliasKey := make([]byte, aliasKeyLen)

	for i := range alliasKey{
		alliasKey[i] = charset[rand.Intn(len(charset))]
	}

	return string(alliasKey)
}

//--------------------------------------------------
//	Get short alias from URL
//	Input:
//		key string - short key
//	Output:
//		alias string - short alias to "longURL"
//func getAlias(key )
