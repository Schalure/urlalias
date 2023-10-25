package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	node, err := h.service.Storage.FindByShortKey(shortKey[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Errorw(
			"error", 
			"err", err.Error(),
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

	if !h.isValidContentType(r, textPlain){
		h.publishBadRequest(&w, fmt.Errorf("content type is not as expected"))
		return
	}

	//	get url
	longURL, err := io.ReadAll(r.Body)
	if err != nil {
		h.publishBadRequest(&w, err)
		h.logger.Info(err.Error())
		return
	}

	if !h.isValidURL(string(longURL)){
		h.publishBadRequest(&w, fmt.Errorf("url is not in the correct format"))
		return
	}

	h.logger.Infow(
		"Parsed URL",
		"Long URL", string(longURL),
	)


	node, err := h.service.Storage.FindByLongURL(string(longURL))
	if err != nil {
		node, err = h.service.NewPairURL(string(longURL)); if err != nil{
			h.publishBadRequest(&w, err)
		}
	}
	aliasURL := h.config.BaseURL() + "/" + node.ShortKey
	h.logger.Infow(
		"Serch/Create alias key",
		"Long URL", node.LongURL,
		"Alias URL", aliasURL,
	)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(aliasURL))
}

func (h *Handlers) checkMainHandlerMethodPost(r *http.Request) error {

	//	execut header "Content-Type" error
	contentType, ok := r.Header["Content-Type"]
	if !ok {
		err := errors.New("header \"Content-Type\" not found")
		h.logger.Info(err.Error())
		return err
	}

	//	execut "Content-Type" value error
	for _, value := range contentType {
		if strings.Contains(value, "text/plain") {
			return nil
		}
	}

	err := fmt.Errorf("error: value of \"content-type\" not right: %s. content-type mast be only \"text/plain\"", contentType)
	h.logger.Info(err.Error())

	return err
}
