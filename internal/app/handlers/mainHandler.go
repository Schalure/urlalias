package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	aliasmaker "github.com/Schalure/urlalias/internal/app/aliasMaker"
	"github.com/Schalure/urlalias/internal/app/storage"
)

// ------------------------------------------------------------
//	"/" request handler.
//	Input:
//		writer http.ResponseWriter
//		request *http.Request
func mainHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		mainHandlerMethodPost(w, r)
	case http.MethodPost:
		mainHandlerMethodGet(w, r)
	default:
		http.Error(w, fmt.Errorf("error: unknown request method: %s", r.Method).Error(), http.StatusBadRequest)
		log.Printf("error: unknown request method: %s\n", r.Method)
	}
}

// ------------------------------------------------------------
//	"/" GET request handler.
//	Execut GET request to make short alias from URL
//	Input:
//		w http.ResponseWriter
//		r *http.Request
func mainHandlerMethodGet(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		longURL, err := storage.GetLongURL(r.RequestURI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		log.Println(longURL)
		w.Header().Add("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(""))
	}
}

// ------------------------------------------------------------
//	"/" POST request handler.
//	Execut POST request to return original URL from short alias
//	Input:
//		w http.ResponseWriter
//		r *http.Request
func mainHandlerMethodPost(w http.ResponseWriter, r *http.Request) {
	if err := checkMainHandlerMethodPost(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	//	get url
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusBadRequest)
		log.Printf("Error: Can't read reqyest body: %s\n", r.Body)
		return
	}

	//	convert data to URL
	log.Println(string(data[:]))

	if aliasURL, err := aliasmaker.GetAliasURL(string(data[:])); err != nil {

		http.Error(w, error.Error(err), http.StatusBadRequest)
		log.Println(err)
		return
	} else {

		log.Println(aliasURL)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(aliasURL))
	}
}

func checkMainHandlerMethodPost(r *http.Request) error {
	/*
		//	execut header "Content-Type" error
		contentType, ok := r.Header["Content-Type"]
		if !ok {
			err := errors.New("header \"Content-Type\" not found")
			log.Println(err.Error())
			return err
		}

		//	execut "Content-Type" value error
		if len(contentType) != 1 || contentType[0] != "text/plain" {
			err := fmt.Errorf("error: value of \"content-type\" not right: %s. content-type mast be only \"text/plain\"", contentType)
			log.Println(err.Error())
			return err
		}
	*/
	return nil
}
