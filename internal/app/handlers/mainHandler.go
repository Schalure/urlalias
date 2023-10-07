package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Schalure/urlalias/cmd/shortener/config"
	aliasmaker "github.com/Schalure/urlalias/internal/app/aliasMaker"
	"github.com/Schalure/urlalias/models"
)

// ------------------------------------------------------------
//	"/" GET request handler.
//	Execut GET request to make short alias from URL
//	Input:
//		w http.ResponseWriter
//		r *http.Request
func MainHandlerMethodGet(repo models.RepositoryURL) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		shortKey := r.RequestURI
		node, err := repo.FindByShortKey(shortKey[1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		log.Println(node.LongURL)
		w.Header().Add("Location", node.LongURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// ------------------------------------------------------------
//	"/" POST request handler.
//	Execut POST request to return original URL from short alias
//	Input:
//		w http.ResponseWriter
//		r *http.Request
func MainHandlerMethodPost(repo models.RepositoryURL) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		if err := checkMainHandlerMethodPost(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}

		//	get url
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(error.Error(err))
			http.Error(w, error.Error(err), http.StatusBadRequest)
			return
		}

		//	Check to valid URL
		u, err := url.ParseRequestURI(string(data[:]))
		if err != nil{
			log.Println(error.Error(err))
			http.Error(w, error.Error(err), http.StatusBadRequest)
			return
		}	
		log.Println(u)


		us := u.String()
		node, err := repo.FindByLongURL(us)
		if err != nil{
			//	try to create alias key
			for i := 0; i < aliasmaker.TrysToMakeAliasKey + 1; i++{
				if i == aliasmaker.TrysToMakeAliasKey{
					log.Println("Can not create alias key")
					http.Error(w, fmt.Errorf("can not create alias key from \"%s\"", u.String()).Error(), http.StatusBadRequest)	
					return		
				}
				aliasKey := aliasmaker.CreateAliasKey()
				node, err = repo.Save(models.AliasURLModel{ID: 0, ShortKey: aliasKey, LongURL: u.String()})
				if err == nil{
					break
				}
			}
		}
		aliasURL := "http://" + *config.ResponseBaseAddres + "/" + node.ShortKey
		log.Printf("Serch/Create alias key: %s - %s\n", node.LongURL, aliasURL)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(aliasURL))		
	}
}


func checkMainHandlerMethodPost(r *http.Request) error {

	//	execut header "Content-Type" error
	contentType, ok := r.Header["Content-Type"]
	if !ok {
		err := errors.New("header \"Content-Type\" not found")
		log.Println(err.Error())
		return err
	}

	//	execut "Content-Type" value error
	for _, value := range contentType{
		if strings.Contains(value, "text/plain"){
			return nil
		}
	}

	err := fmt.Errorf("error: value of \"content-type\" not right: %s. content-type mast be only \"text/plain\"", contentType)
	log.Println(err.Error())
	return err
}
