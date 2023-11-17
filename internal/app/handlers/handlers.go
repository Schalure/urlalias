package handlers

import (
	"net/url"

	"github.com/Schalure/urlalias/internal/app/aliasmaker"
)

const (
	contentType     string = "Content-Type"
	contentEncoding string = "Content-Encoding"
	acceptEncoding  string = "Accept-Encoding"
)

const (
	textPlain = "text/plain"
	appJSON   = "application/json"
)

var ContentTypeToCompress = []string{
	textPlain,
	appJSON,
}

type Handlers struct {
	service *aliasmaker.AliasMakerServise
}

// ------------------------------------------------------------
//
//	Constructor of Handlers type
func NewHandlers(service *aliasmaker.AliasMakerServise) *Handlers {

	return &Handlers{
		service: service,
	}
}

// ------------------------------------------------------------
//
//	Check to valid URL - method of Handlers type
func (h *Handlers) isValidURL(u string) bool {

	if _, err := url.ParseRequestURI(u); err != nil {
		h.service.Logger.Infow(
			"URL is not in the correct format",
			"URL", u,
		)
		return false
	}
	return true
}
