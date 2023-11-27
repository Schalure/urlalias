package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Schalure/urlalias/internal/app/interpreter"
	"github.com/Schalure/urlalias/internal/app/models"
)

// ------------------------------------------------------------
//	POST request - "/"
func (h *Handlers) APIShortenHandlerPost(w http.ResponseWriter, r *http.Request) {

	type (
		requestModel struct {
			URL string `json:"url"`
		}
		responseModel struct {
			Result string `json:"result"`
		}
	)

	var (
		requestJSON requestModel
		i           interpreter.InterpreterJSON
	)

	userID := r.Context().Value(UserID)
	uID, ok := userID.(uint64)
	if !ok {
		http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
		h.service.Logger.Info(errors.New("can't parsed user id").Error())
		return
	}

	err := i.Unmarshal(r.Body, &requestJSON)
	if err != nil {
		http.Error(w, "can't decode JSON content", http.StatusBadRequest)
		h.service.Logger.Infow(
			"Can't decode JSON content",
			"err", err.Error(),
		)
		return
	}

	if !h.isValidURL(requestJSON.URL) {
		http.Error(w, "url is not in the correct format", http.StatusBadRequest)
		h.service.Logger.Infow(
			"url is not in the correct format",
			"url", requestJSON.URL,
		)
		return
	}


	node, statusCode, err := h.service.CreateAlias(uID, string(requestJSON.URL))
	if err != nil {
		h.service.Logger.Infow("alias to save", "error", err)	
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.service.Logger.Infow("alias to save",
		"status", statusCode,
		"node.ID", node.ID,
		"node.UserID", node.UserID,
		"node.LongURL", node.LongURL,
		"node.ShortKey", node.ShortKey,
	)

	var resp = responseModel{
		Result: h.service.Config.BaseURL() + "/" + node.ShortKey,
	}
	buf, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.service.Logger.Infow(
			"Can not encode data",
			"data", resp,
			"err", err,
		)
		return
	}

	w.Header().Set("Content-Type", appJSON)
	w.WriteHeader(statusCode)
	w.Write(buf)
}

func (h *Handlers) APIShortenBatchHandlerPost(w http.ResponseWriter, r *http.Request) {

	type (
		requestModel struct {
			ID          string `json:"correlation_id"`
			OriginalURL string `json:"original_url"`
		}

		responseModel struct {
			ID       string `json:"correlation_id"`
			ShortURL string `json:"short_url"`
		}
	)

	var (
		requestJSON  []requestModel
		responseJSON []responseModel
		i            interpreter.InterpreterJSON
		nodes        []models.AliasURLModel
	)

	userID := r.Context().Value(UserID)
	uID, ok := userID.(uint64)
	if !ok {
		http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
		h.service.Logger.Info(errors.New("can't parsed user id").Error())
		return
	}

	err := i.Unmarshal(r.Body, &requestJSON)
	if err != nil {
		http.Error(w, fmt.Sprintf("can't decode JSON content, error: %s", err), http.StatusBadRequest)
		h.service.Logger.Infow(
			"Can't decode JSON content",
			"err", err.Error(),
		)
		return
	}

	for _, req := range requestJSON {

		node := h.service.Storage.FindByLongURL(req.OriginalURL)
		if node == nil {
			node, err = h.service.NewPairURL(req.OriginalURL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				h.service.Logger.Infow(
					"Can't create pair url",
					"err", err.Error(),
				)
				return
			}
			node.UserID = uID
		}
		nodes = append(nodes, *node)
		responseJSON = append(responseJSON, responseModel{req.ID, h.service.Config.BaseURL() + "/" + node.ShortKey})
	}

	if err := h.service.Storage.SaveAll(nodes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.service.Logger.Infow(
			"Can't save to storage",
			"err", err.Error(),
		)
		return
	}

	buf, err := json.Marshal(&responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.service.Logger.Infow(
			"Can't dekode to JSON",
			"buf", string(buf),
			"err", err.Error(),
		)
		return
	}

	w.Header().Set("Content-Type", appJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}

func (h *Handlers) APIUserURLsHandlerGet(w http.ResponseWriter, r *http.Request) {

	type responseModel struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	var (
		responseJSON []responseModel
	)

	userID := r.Context().Value(UserID)
	uID, ok := userID.(uint64)
	if !ok {
		http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
		h.service.Logger.Info(errors.New("can't parsed user id").Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	nodes, err := h.service.Storage.FindByUserID(ctx, uID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(nodes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for _, node := range nodes {
		responseNodeJSON := responseModel{
			ShortURL:    h.service.Config.BaseURL() + "/" + node.ShortKey,
			OriginalURL: node.LongURL,
		}
		responseJSON = append(responseJSON, responseNodeJSON)
	}

	buf, err := json.Marshal(&responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.service.Logger.Infow(
			"Can't dekode to JSON",
			"buf", string(buf),
			"err", err.Error(),
		)
		return
	}

	w.Header().Set("Content-Type", appJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func (h *Handlers) APIUserURLsHandlerDelete(w http.ResponseWriter, r *http.Request) {

	var(
		aliases []string
		i interpreter.InterpreterJSON
	)
	userID := r.Context().Value(UserID)
	uID, ok := userID.(uint64)
	if !ok {
		http.Error(w, errors.New("can't parsed user id").Error(), http.StatusBadRequest)
		h.service.Logger.Info(errors.New("can't parsed user id").Error())
		return
	}

	if err := i.Unmarshal(r.Body, &aliases); err != nil {
		http.Error(w, fmt.Sprintf("can't decode JSON content, error: %s", err), http.StatusBadRequest)
		h.service.Logger.Infow(
			"Can't decode JSON content",
			"err", err.Error(),
		)
		return
	}

	go h.service.DeleteUserURLs(uID, aliases)

	w.WriteHeader(http.StatusAccepted)
}

//	Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDEwMzkyOTQsIlVzZXJJRCI6MzB9.w8j0xOKSrgLwTg7_tESoscCcmIx2SBTSW0ftwtoft8g
