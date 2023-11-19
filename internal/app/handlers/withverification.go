package handlers

import (
	"context"
	"errors"
	"net/http"
)

func (m *Middleware) WithVerification(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var userID uint64

		tokenCookie, err := r.Cookie(authorization)
		if err != nil {
			http.Error(w, errors.New("Unauthorized").Error(), http.StatusUnauthorized)
			return
		} else if userID, err = getUserID(tokenCookie.Value); err != nil {
			http.Error(w, errors.New("Unauthorized").Error(), http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userID", userID)))
	})
}
