package handlers

import "net/http"

func (m *Middleware) WithAuthentication(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}