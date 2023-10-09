package storage

import (
	"errors"
	"fmt"
)

// Map to save pair
//
//	key - "shortURL
//	val - longURL"
var dataBase = map[string]string{}

// --------------------------------------------------
//
//	Get alias from DB
//	Input:
//		longURL string - valid URL
//	Output:
//		urlInfo URLInfo - short alias, and other info about "longURL"
//		ok bool - true: alias was found, false - alias was not found
func GetAliasKey(longURL string) (string, bool) {

	for k, v := range dataBase {
		if v == longURL {
			return k, true
		}
	}
	return "", false
}

func GetLongURL(aliasKey string) (string, error) {

	if longURL, ok := dataBase[aliasKey]; !ok {
		return "", fmt.Errorf("can't find the url by key: %s", aliasKey)
	} else {
		return longURL, nil
	}
}

// --------------------------------------------------
//
//	Save pair "longURL, alias" to DB
//	Input:
//		longURL string - valid URL
//		alias string - alias to valid URL
//	Output:
//		err error - can not save the repeated value of short url
func SavePair(longURL, shortURL string) error {

	if _, ok := dataBase[shortURL]; !ok {
		dataBase[shortURL] = longURL
		return nil
	}
	return errors.New("can't save the repeated value of short url")

}
