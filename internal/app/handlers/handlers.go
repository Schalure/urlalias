package handlers

import (
	"net/url"

	"github.com/Schalure/urlalias/cmd/shortener/config"
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
	config  *config.Configuration
	logger  Loggerer
}

// ------------------------------------------------------------
//
//	Constructor of Handlers type
//	Input:
//		storage IStorage
//		config *config.Configuration
//		logger *zap.SugaredLogger
//	Output:
//		*Handlers - ptr to new Handlers
func NewHandlers(service *aliasmaker.AliasMakerServise, config *config.Configuration, logger Loggerer) *Handlers {

	return &Handlers{
		service: service,
		config:  config,
		logger:  logger,
	}
}

// ------------------------------------------------------------
//
//	Check to valid URL - method of Handlers type
//	Receiver:
//		h* Handlers
//	Input:
//		url string
//	Output:
//		bool
func (h *Handlers) isValidURL(u string) bool {

	if _, err := url.ParseRequestURI(u); err != nil {
		h.logger.Infow(
			"URL is not in the correct format",
			"URL", u,
		)
		return false
	}
	return true
}
