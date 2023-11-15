package handlers

import (
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
		h.logger.Infow(
			"The urlAliasNode not found by key",
			"Key", shortKey,
		)
		return
	}
	h.logger.Infow(
		"Long URL",
		"URL", node.LongURL,
	)

	w.Header().Add("Location", node.LongURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Info(err.Error())
		return
	}

	if !h.isValidURL(string(longURL)) {
		http.Error(w, fmt.Sprintf("url is not in the correct format: %s", longURL), http.StatusBadRequest)
		return
	}

	h.logger.Infow(
		"Parsed URL",
		"Long URL", string(longURL),
	)

	var statusCode int
	node := h.service.Storage.FindByLongURL(string(longURL))
	if node == nil {
		if node, err = h.service.NewPairURL(string(longURL)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.logger.Info(err.Error())
			return
		}
		if err = h.service.Storage.Save(node); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.logger.Info(err.Error())
			return
		}
		statusCode = http.StatusCreated
	} else {
		statusCode = http.StatusConflict
	}

	aliasURL := h.config.BaseURL() + "/" + node.ShortKey

	h.logger.Infow(
		"Serch/Create alias key",
		"Long URL", node.LongURL,
		"Alias URL", aliasURL,
	)

	w.Header().Set("Content-Type", textPlain)
	w.WriteHeader(statusCode)
	w.Write([]byte(aliasURL))
}
