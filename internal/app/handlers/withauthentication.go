package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ContextKey string

const UserID ContextKey = "userID"

const tokenExp = time.Hour * 3
const secretKey = "supersecretkey"

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
			m.service.Logger.Infow(
				"WithAuthentication: tokenCookie, err := r.Cookie(authorization)",
				"error", err,
			)
			if userID, err = m.service.CreateUser(); err != nil {
				m.service.Logger.Infow(
					"WithAuthentication: userID, err = m.service.CreateUser()",
					"error", err,
				)
				http.Error(w, errors.New("internal error").Error(), http.StatusInternalServerError)
				return
			}
			tokenString, err = createTokenJWT(userID)
			if err != nil {
				m.service.Logger.Infow(
					"WithAuthentication: tokenString, err = createTokenJWT(userID)",
					"error", err,
				)
				http.Error(w, errors.New("internal error").Error(), http.StatusInternalServerError)
				return
			}

			m.service.Logger.Infow(
				"Add new user",
				"userID", userID,
			)
			http.SetCookie(w, &http.Cookie{
				Name:  authorization,
				Value: tokenString,
			})

		} else if userID, err = getUserID(tokenCookie.Value); err != nil {
			m.service.Logger.Infow(
				"WithAuthentication: userID, err = getUserID(tokenCookie.Value)",
				"error", err,
			)
			if userID, err = m.service.CreateUser(); err != nil {
				m.service.Logger.Infow(
					"WithAuthentication: userID, err = m.service.CreateUser()",
					"error", err,
				)
				http.Error(w, errors.New("internal error").Error(), http.StatusInternalServerError)
				return
			}
			tokenString, err = createTokenJWT(userID)
			if err != nil {
				m.service.Logger.Infow(
					"WithAuthentication: tokenString, err = createTokenJWT(userID)",
					"error", err,
				)
				http.Error(w, errors.New("internal error").Error(), http.StatusInternalServerError)
				return
			}

			m.service.Logger.Infow(
				"Add new user",
				"userID", userID,
			)
			http.SetCookie(w, &http.Cookie{
				Name:  authorization,
				Value: tokenString,
			})
		} else {
			authCookie, _ := r.Cookie(authorization)
			http.SetCookie(w, authCookie)
		}

		m.service.Logger.Infow(
			"Request from user",
			"userID", userID,
		)

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserID, userID)))
	})
}

// get user id from JWT token string
func getUserID(tokenString string) (uint64, error) {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
