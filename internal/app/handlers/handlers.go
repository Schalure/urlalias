package handlers

import (
	"net/http"
)

// --------------------------------------------------
var (
	//	Hadler func list
	HandlersList = map[string]http.HandlerFunc{
		"/": mainHandler,
	}
)

