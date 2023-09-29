package app

import "net/http"


var(
	//	Hadler func list
	HandlersList = map[string]http.HandlerFunc {
		"/" : mainHandler,
	}
)

//--------------------------------------------------
//	"/" request handler
func mainHandler(writer http.ResponseWriter, request *http.Request){
	writer.Write([]byte("hi"))
}