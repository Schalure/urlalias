package app

import (
	"io"
	"net/http"
)

//--------------------------------------------------
var(
	//	Hadler func list
	HandlersList = map[string]http.HandlerFunc {
		"/" : mainHandler,
	}
)

//--------------------------------------------------
//	"/" request handler
func mainHandler(writer http.ResponseWriter, request *http.Request){
	
	//	only POST request to execut
	if request.Method != http.MethodPost{
		http.Error(writer, "only POST requests are accepted on the path \"/\"", http.StatusBadRequest)
		return
	}

	//	execut header "Content-Type" error
	contentType, ok := request.Header["Content-Type"]; 
	if !ok{
		http.Error(writer, "header \"Content-Type\" not found", http.StatusBadRequest)
		return
	}

	//	execut "Content-Type" value error
	if len(contentType) != 1 || contentType[0] != "text/plain"{
		http.Error(writer, "Content-Type mast be only \"text/plain\"", http.StatusBadRequest)
		return
	}

	//	get url
	data, err := io.ReadAll(request.Body)
	if err != nil{
		http.Error(writer, error.Error(err), http.StatusBadRequest)
		return
	}

	//	convert data to URL
	url := string(data[:])
	if shortURL, err := makeAliasUrl(url); err != nil{
		http.Error(writer, error.Error(err), http.StatusBadRequest)
		return
	}else{
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte(shortURL))
	}

}