package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Schalure/urlalias/internal/app/interpreter"
	"github.com/Schalure/urlalias/internal/app/storage"
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
		h.publishBadRequest(&w, fmt.Errorf("can't decode JSON content"))
		h.logger.Infow(
			"Can't decode JSON content",
			"err", err.Error(),
		)
		return
	}

	if !h.isValidURL(requestJSON.URL) {
		h.publishBadRequest(&w, fmt.Errorf("url is not in the correct format"))
		return
	}

	node := h.service.Storage.FindByLongURL(string(requestJSON.URL))
	if node == nil {
		if node, err = h.service.NewPairURL(string(requestJSON.URL)); err != nil {
			h.publishBadRequest(&w, err)
			return
		}
		if err = h.service.Storage.Save(node); err != nil{
			h.publishBadRequest(&w, err)
			return
		}
	}

	var resp = responseModel{
		Result: h.config.BaseURL() + "/" + node.ShortKey,
	}
	buf, err := json.Marshal(&resp)
	if err != nil {
		h.publishBadRequest(&w, err)
		h.logger.Infow(
			"Can not encode data",
			"data", resp,
			"err", err,
		)
		return
	}

	h.logger.Infow(
		"Serch/Create alias key",
		"Long URL", node.LongURL,
		"Alias URL", resp.Result,
	)

	w.Header().Set("Content-Type", appJSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(buf)
}

// /api/shorten/batch
// request:
// [
//
//	{
//	    "correlation_id": "<строковый идентификатор>",
//	    "original_url": "<URL для сокращения>"
//	},
//	...
//
// ]
// response:
// [
//
//	{
//	    "correlation_id": "<строковый идентификатор из объекта запроса>",
//	    "short_url": "<результирующий сокращённый URL>"
//	},
//	...
//
// ]
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
		nodes        []storage.AliasURLModel
	)

	err := i.Unmarshal(r.Body, &requestJSON)
	if err != nil {
		h.publishBadRequest(&w, fmt.Errorf("can't decode JSON content"))
		h.logger.Infow(
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
				h.publishBadRequest(&w, fmt.Errorf("can't decode JSON content"))
				h.logger.Infow(
					"Can't create pair url",
					"err", err.Error(),
				)
				return
			}
		}
		nodes = append(nodes, *node)
		responseJSON = append(responseJSON, responseModel{req.ID, h.config.BaseURL() + "/" + node.ShortKey})
	}

	if err := h.service.Storage.SaveAll(nodes); err != nil {
		h.publishBadRequest(&w, fmt.Errorf("can't decode JSON content"))
		h.logger.Infow(
			"Can't save to storage",
			"err", err.Error(),
		)
		return
	}

	buf, err := json.Marshal(&responseJSON)
	if err != nil {
		h.publishBadRequest(&w, fmt.Errorf("can't decode JSON content"))
		h.logger.Infow(
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
