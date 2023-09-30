package app

import (
	"io"
	"log"
	"net/http"

	aliasmaker "github.com/Schalure/urlalias/internal/app/aliasMaker"
)

//--------------------------------------------------
var(
	//	Hadler func list
	HandlersList = map[string]http.HandlerFunc {
		"/" : mainHandler,
	}
)

//--------------------------------------------------
//	"/" request handler.
//	Execut request to make short alias from URL
//	Input:
//		writer http.ResponseWriter
//		request *http.Request
func mainHandler(writer http.ResponseWriter, request *http.Request){
	
	//	only POST request to execut
	if request.Method != http.MethodPost{
		http.Error(writer, "only POST requests are accepted on the path \"/\"", http.StatusBadRequest)
		log.Printf("Error: method %s only POST requests are accepted on the path \"/\"\n", request.Method)
		return
	}

	//	execut header "Content-Type" error
	contentType, ok := request.Header["Content-Type"]; 
	if !ok{
		http.Error(writer, "header \"Content-Type\" not found", http.StatusBadRequest)
		log.Printf("Error: header \"Content-Type\" not found\n")
		return
	}

	//	execut "Content-Type" value error
	if len(contentType) != 1 || contentType[0] != "text/plain"{
		http.Error(writer, "Content-Type mast be only \"text/plain\"", http.StatusBadRequest)
		log.Printf("Error: value of \"Content-Type\" not right: %s. Content-Type mast be only \"text/plain\"\n", contentType)
		return
	}

	//	get url
	data, err := io.ReadAll(request.Body)
	if err != nil{
		http.Error(writer, error.Error(err), http.StatusBadRequest)
		log.Printf("Error: Can't read reqyest body: %s\n", request.Body)
		return
	}

	//	convert data to URL
	url := string(data[:])
	log.Println(url)

	if shortURL, err := aliasmaker.MakeAliasUrl(url); err != nil{
		http.Error(writer, error.Error(err), http.StatusBadRequest)
		log.Println(err)
		return
	}else{
		log.Println(shortURL)
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte(shortURL))
	}
}