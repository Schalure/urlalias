package handlers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Schalure/urlalias/internal/app/aliasmaker"
)

const (
	contentType     string = "Content-Type"
	contentEncoding string = "Content-Encoding"
	acceptEncoding  string = "Accept-Encoding"
	authorization   string = "Authorization"
)

const (
	textPlain = "text/plain"
	appJSON   = "application/json"
)

var ContentTypeToCompress = []string{
	textPlain,
	appJSON,
}

//go:generate mockgen -destination=../mocks/mock_shortner.go -package=mocks github.com/Schalure/gofermart/internal/handlers Shortner
type Shortner interface {
	GetOriginalURL(ctx context.Context, shortKey string) (string, error)
}



type Handler struct {
	service *aliasmaker.AliasMakerServise
}

// ------------------------------------------------------------
//
//	Constructor of Handlers type
func NewHandlers(service *aliasmaker.AliasMakerServise) *Handler {

	return &Handler{
		service: service,
	}
}

// ------------------------------------------------------------
//
//	Check to valid URL - method of Handlers type
func (h *Handler) isValidURL(u string) bool {

	if _, err := url.ParseRequestURI(u); err != nil {
		h.service.Logger.Infow(
			"URL is not in the correct format",
			"URL", u,
		)
		return false
	}
	return true
}

// Get login from request context
func (h *Handler) getLoginFromContext(ctx context.Context) (string, error) {

	login := ctx.Value(UserID)
	l, ok := login.(string)
	if !ok {
		return "", fmt.Errorf("login is not valid")
	}
	return l, nil
}