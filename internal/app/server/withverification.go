package server

import (
	"context"
	"errors"
	"net/http"
)

func (m *Middleware) WithVerification(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenCookie, err := r.Cookie(authorization)
		if err != nil {
			m.logger.Infow(
				"WithVerification: tokenCookie, err := r.Cookie(authorization)",
				"error", err,
			)
			http.Error(w, errors.New("Unauthorized").Error(), http.StatusUnauthorized)
			return
		}

		userID, err := getUserID(tokenCookie.Value)
		if err != nil {
			m.logger.Infow(
				"WithVerification: userID, err = getUserID(tokenCookie.Value)",
				"error", err,
				"user", userID,
				"token", tokenCookie.Value,
			)
			http.Error(w, errors.New("Unauthorized").Error(), http.StatusUnauthorized)
			return
		}

		authCookie, _ := r.Cookie(authorization)
		http.SetCookie(w, authCookie)

		m.logger.Infow(
			"Request from user",
			"userID", userID,
		)

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserID, userID)))
	})
}
