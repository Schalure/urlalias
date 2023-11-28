package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ------------------------------------------------------------
//
//	"/" GET request handler.
//	Execut GET request to make short alias from URL
//	Input:
//		w http.ResponseWriter
//		r *http.Request
func (h *Handlers) mainHandlerGet(w http.ResponseWriter, r *http.Request) {

	shortKey := r.RequestURI
	node := h.service.Storage.FindByShortKey(shortKey[1:])
	if node == nil {
		http.Error(w, fmt.Sprintf("the urlAliasNode not found by key \"%s\"", shortKey), http.StatusBadRequest)
		h.service.Logger.Infow(
			"The urlAliasNode not found by key",
			"Key", shortKey,
		)
		return
	}
	h.service.Logger.Infow(
		"Long URL",
		"URL", node.LongURL,
	)

	var status int
	if node.DeletedFlag {
		status = http.StatusGone
	} else {
		w.Header().Add("Location", node.LongURL)
		status = http.StatusTemporaryRedirect
	}

	w.WriteHeader(status)
}

// ------------------------------------------------------------
//
//	"/" POST request handler.
//	Execut POST request to return original URL from short alias
//	Input:
//		w http.ResponseWriter
//		r *http.Request
func (h *Handlers) mainHandlerPost(w http.ResponseWriter, r *http.Request) {

	//	get url
	longURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("can`t read request body: %s", err.Error()).Error(), http.StatusBadRequest)
		h.service.Logger.Info(err.Error())
		return
	}

	if !h.isValidURL(string(longURL)) {
		http.Error(w, fmt.Sprintf("url is not in the correct format: %s", longURL), http.StatusBadRequest)
		return
	}

	h.service.Logger.Infow(
		"Parsed URL",
		"Long URL", string(longURL),
	)

	var statusCode int
	node := h.service.Storage.FindByLongURL(string(longURL))
	if node == nil {
		if node, err = h.service.NewPairURL(string(longURL)); err != nil {
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
