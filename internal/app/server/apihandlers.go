package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/interpreter"
)

//	Handler retuns short URL by original URL. Handler can returns three HTTP statuses:
//	1. StatusBadRequest (400) - if an internal service error occurred;
//	2. StatusConflict (409) - if the original URL is already saved in the service;
//	3. StatusCreated (201) - if original URL is saved successfully and alias is created.
func (h *Server) apiGetShortURL(w http.ResponseWriter, r *http.Request) {

	type (
		RequestJSON struct {
			OriginalURL string `json:"url"`
		}
		ResponseJSON struct {
			ShortURL string `json:"result"`
		}
	)

	var i interpreter.InterpreterJSON

	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
		return
	}

	var requestJSON RequestJSON
	err = i.Unmarshal(r.Body, &requestJSON)
	if err != nil {
		http.Error(w, "can't decode JSON content", http.StatusBadRequest)
		return
	}

	var statusCode int
	shortURL, err := h.shortner.GetShortKey(r.Context(), userID, requestJSON.OriginalURL)
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

	buf, err := json.Marshal(&ResponseJSON{ShortURL: h.baseURL + "/" + shortURL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", appJSON)
	w.WriteHeader(statusCode)
	w.Write(buf)
}


func (h *Server) apiGetBatchShortURL(w http.ResponseWriter, r *http.Request) {

	type (
		RequestJSON struct {
			ID          string `json:"correlation_id"`
			OriginalURL string `json:"original_url"`
		}

		ResponseJSON struct {
			ID       string `json:"correlation_id"`
			ShortURL string `json:"short_url"`
		}
	)

	var (
		requestJSON  []RequestJSON
		responseJSON []ResponseJSON
		i            interpreter.InterpreterJSON
	)

	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
		return
	}

	err = i.Unmarshal(r.Body, &requestJSON)
	if err != nil {
		http.Error(w, fmt.Sprintf("can't decode JSON content, error: %s", err), http.StatusBadRequest)
		return
	}

	batchOriginalURL := make([]string, len(requestJSON))
	for i, request := range requestJSON {
		batchOriginalURL[i] = request.OriginalURL
	}

	batchShortKey, err := h.shortner.GetBatchShortURL(r.Context(), userID, batchOriginalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Infow("Can't save to storage", "err", err.Error(),)
		return		
	}

	responseJSON = make([]ResponseJSON, len(requestJSON))
	for i, shortKey := range batchShortKey {
		responseJSON[i].ID = requestJSON[i].ID
		responseJSON[i].ShortURL = h.baseURL + "/" + shortKey
	}


	buf, err := json.Marshal(&responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Infow("Can't dekode to JSON", "buf", string(buf), "err", err.Error())
		return
	}

	w.Header().Set("Content-Type", appJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}


func (h *Server) apiGetUserAliases(w http.ResponseWriter, r *http.Request) {

	type responseModel struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	var responseJSON []responseModel

	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
		return
	}


	ctx, cancel := context.WithTimeout(r.Context(), time.Second * 1)
	defer cancel()

	nodes, err := h.userManager.GetUserAliases(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(nodes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for _, node := range nodes {
		responseNodeJSON := responseModel{
			ShortURL:    h.baseURL + "/" + node.ShortKey,
			OriginalURL: node.LongURL,
		}
		responseJSON = append(responseJSON, responseNodeJSON)
	}

	buf, err := json.Marshal(&responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", appJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}


func (h *Server) aipDeleteUserAliases(w http.ResponseWriter, r *http.Request) {

	var (
		aliases []string
		i       interpreter.InterpreterJSON
	)

	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
		return
	}

	if err := i.Unmarshal(r.Body, &aliases); err != nil {
		http.Error(w, fmt.Sprintf("can't decode JSON content, error: %s", err), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second * 5)
	defer cancel()
	if err := h.shortner.AddAliasesToDelete(ctx, userID, aliases...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

//	Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDEwMzkyOTQsIlVzZXJJRCI6MzB9.w8j0xOKSrgLwTg7_tESoscCcmIx2SBTSW0ftwtoft8g
