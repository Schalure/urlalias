package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Schalure/urlalias/internal/app/interpreter"
	"github.com/Schalure/urlalias/internal/app/models"
)

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

	var statusCode int
	node := h.service.Storage.FindByLongURL(string(requestJSON.URL))
	if node == nil {
		if node, err = h.service.NewPairURL(string(requestJSON.URL)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.service.Logger.Info(err.Error())
			return
		}
		if err = h.service.Storage.Save(node); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.service.Logger.Info(err.Error())
			return
		}
		statusCode = http.StatusCreated
	} else {
		statusCode = http.StatusConflict
	}

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

	h.service.Logger.Infow(
		"Serch/Create alias key",
		"Long URL", node.LongURL,
		"Alias URL", resp.Result,
	)

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
	
}
