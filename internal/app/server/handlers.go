package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/models/aliasentity"
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

//go:generate mockgen -destination=../mocks/mock_shortner.go -package=mocks github.com/Schalure/urlalias/internal/app/handlers Shortner
type Shortner interface {
	GetOriginalURL(ctx context.Context, shortKey string) (string, error)
	GetShortKey(ctx context.Context, userID uint64, originalURL string) (string, error)
	GetBatchShortURL(ctx context.Context, userID uint64, batchOriginalURL []string) ([]string, error)
	AddAliasesToDelete(ctx context.Context, userID uint64, aliases ...string) error
	IsDatabaseActive() bool
}

//go:generate mockgen -destination=../mocks/mock_usermanager.go -package=mocks github.com/Schalure/urlalias/internal/app/handlers UserManager
type UserManager interface {
	CreateUser() (uint64, error)
	GetUserAliases(ctx context.Context, userID uint64) ([]aliasentity.AliasURLModel, error)
}


type Server struct {
	userManager UserManager
	shortner Shortner

	logger *zaplogger.ZapLogger

	baseURL string
}


//	Constructor of Handler type
func New(userManager UserManager, shortner Shortner, logger *zaplogger.ZapLogger, baseURL string) *Server {

	return &Server{
		userManager: userManager,
		shortner: shortner,
		logger: logger,
		baseURL: baseURL,
	}
}


//	Handler retuns original URL by short key in HTTP header "Location" and redirect status code (307).
//	If URL not found or was deleted, returns error
func (h *Server) redirect(w http.ResponseWriter, r *http.Request) {

	shortKey := r.RequestURI[1:]

	originalURL, err := h.shortner.GetOriginalURL(r.Context(), shortKey)
	if err != nil {
		if errors.Is(err, aliasmaker.ErrURLNotFound){
			http.Error(w, fmt.Sprintf("the url alias not found by key \"%s\"", shortKey), http.StatusBadRequest)
			return
		}
		if errors.Is(err, aliasmaker.ErrURLWasDeleted){
			http.Error(w, fmt.Sprintf("the url alias was deleted \"%s\"", shortKey), http.StatusGone)
			return		
		}
	}

	w.Header().Add("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}


//	Handler retuns short URL by original URL. Handler can returns three HTTP statuses:
//	1. StatusBadRequest (400) - if an internal service error occurred;
//	2. StatusConflict (409) - if the original URL is already saved in the service;
//	3. StatusCreated (201) - if original URL is saved successfully and alias is created.
func (h *Server) getShortURL(w http.ResponseWriter, r *http.Request) {

	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
		return
	}

	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("can`t read request body: %s", err.Error()).Error(), http.StatusBadRequest)
		return
	}


	var statusCode int
	shortURL, err := h.shortner.GetShortKey(r.Context(), userID, string(originalURL))
	if err != nil {
		if errors.Is(err, aliasmaker.ErrInternal) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, aliasmaker.ErrConflictURL) {
			statusCode = http.StatusConflict
		}
	} else {
		statusCode = http.StatusCreated
	}

	w.Header().Set("Content-Type", textPlain)
	w.WriteHeader(statusCode)
	w.Write([]byte(h.baseURL + "/" + shortURL))
}


//	Get state of database service
func (h *Server) PingGet(w http.ResponseWriter, r *http.Request) {

	if !h.shortner.IsDatabaseActive() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}


// Get User ID from request context
func (h *Server) getUserIDFromContext(ctx context.Context) (uint64, error) {

	userID := ctx.Value(UserID)
	ID, ok := userID.(uint64)
	if !ok {
		return 0, fmt.Errorf("login is not valid")
	}
	return ID, nil
}