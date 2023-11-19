package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

func (m *Middleware) WithAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var tokenString string
		var userID uint64

		tokenCookie, err := r.Cookie(authorization)
		if err != nil {

			if userID, err = m.service.CreateUser(); err != nil {
				http.Error(w, errors.New("internal error").Error(), http.StatusInternalServerError)
				return
			}
			tokenString, err = createTokenJWT(userID)
			if err != nil {
				http.Error(w, errors.New("internal error").Error(), http.StatusInternalServerError)
				return
			}
		} else if userID, err = getUserID(tokenCookie.Value); err != nil {
			if userID, err = m.service.CreateUser(); err != nil {
				http.Error(w, errors.New("internal error").Error(), http.StatusInternalServerError)
				return
			}
			tokenString, err = createTokenJWT(userID)
			if err != nil {
				http.Error(w, errors.New("internal error").Error(), http.StatusInternalServerError)
				return
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:  authorization,
			Value: tokenString,
		})
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userID", userID)))
	})
}

// get user id from JWT token string
func getUserID(tokenString string) (uint64, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return 0, errors.New("can't parse token string")
	}
	if !token.Valid {
		return 0, errors.New("token not valid")
	}
	return claims.UserID, nil
}

func createTokenJWT(userID uint64) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
