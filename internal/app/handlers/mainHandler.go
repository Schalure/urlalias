package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Schalure/urlalias/internal/app/aliasmaker"
)

//	Handler retuns original URL by short key in HTTP header "Location" and redirect status code (307).
//	If URL not found or was deleted, returns error
func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {

	shortKey := r.RequestURI[1:]

	originalURL, err := h.service.GetOriginalURL(r.Context(), shortKey)
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

// ------------------------------------------------------------
//
//	"/" POST request handler.
//	Execut POST request to return original URL from short alias
//	Input:
//		w http.ResponseWriter
//		r *http.Request
func (h *Handler) mainHandlerPost(w http.ResponseWriter, r *http.Request) {

	//	get url
	originalURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("can`t read request body: %s", err.Error()).Error(), http.StatusBadRequest)
		h.service.Logger.Info(err.Error())
		return
	}

	if !h.isValidURL(string(originalURL)) {
		http.Error(w, fmt.Sprintf("url is not in the correct format: %s", originalURL), http.StatusBadRequest)
		return
	}

	h.service.Logger.Infow(
		"Parsed URL",
		"Long URL", string(originalURL),
	)

	


	var statusCode int
	node := h.service.Storage.FindByLongURL(string(originalURL))
	if node == nil {
		if node, err = h.service.NewPairURL(string(originalURL)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.service.Logger.Info(err.Error())
			return
		}

		userID := r.Context().Value(UserID)
		uID, ok := userID.(uint64)
		if !ok {
			http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
			h.service.Logger.Info(errors.New("can't parsed user id").Error())
			return
		}
		node.UserID = uID

		if err = h.service.Storage.Save(node); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.service.Logger.Info(err.Error())
			return
		}
		statusCode = http.StatusCreated
	} else {
		statusCode = http.StatusConflict
	}

	aliasURL := h.service.Config.BaseURL() + "/" + node.ShortKey

	h.service.Logger.Infow(
		"Serch/Create alias key",
		"Long URL", node.LongURL,
		"Alias URL", aliasURL,
	)

	w.Header().Set("Content-Type", textPlain)
	w.WriteHeader(statusCode)
	w.Write([]byte(aliasURL))
}
