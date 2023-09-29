package app

import (
	"net/url"
)

//--------------------------------------------------
//	Make short alias from URL
func makeAliasUrl(longURL string) (string, error){

	//	Check to valid URL
	if _, err := url.ParseRequestURI(longURL); err != nil{
		return "", err
	}

	alias, ok := getAliasFromDB(longURL)
	if !ok{
		alias = (longURL)
		if err := savePairToDB(longURL, alias); err != nil{
			
		}
	}

	return alias, nil
}
