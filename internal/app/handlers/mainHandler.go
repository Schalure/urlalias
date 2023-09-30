package handlers

import (
	"io"
	"log"
	"net/http"

	aliasmaker "github.com/Schalure/urlalias/internal/app/aliasMaker"
	"github.com/Schalure/urlalias/internal/app/database"
)

// --------------------------------------------------
//
//	"/" request handler.
//	Execut request to make short alias from URL
//	Input:
//		writer http.ResponseWriter
//		request *http.Request
func mainHandler(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodGet {
		longURL, err := database.GetLongURL(request.RequestURI)
		if err != nil{
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		log.Println(longURL)
		writer.Header().Add("Location", longURL)
		writer.WriteHeader(http.StatusTemporaryRedirect)
		writer.Write([]byte(""))
	}

	//	only POST request to execut
	if request.Method == http.MethodPost {

		if err := checkMainHandlerMethodPost(request); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}

		//	get url
		data, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, error.Error(err), http.StatusBadRequest)
			log.Printf("Error: Can't read reqyest body: %s\n", request.Body)
			return
		}

		//	convert data to URL
		log.Println(string(data[:]))

		if aliasURL, err := aliasmaker.GetAliasURL(string(data[:])); err != nil {

			http.Error(writer, error.Error(err), http.StatusBadRequest)
			log.Println(err)
			return
		} else {

			log.Println(aliasURL)
			writer.Header().Set("Content-Type", "text/plain")
			writer.WriteHeader(http.StatusCreated)
			writer.Write([]byte(aliasURL))
		}

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

