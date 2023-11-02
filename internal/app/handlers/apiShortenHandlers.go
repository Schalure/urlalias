package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Schalure/urlalias/internal/app/interpreter"
)

type requestModel struct {
	URL string `json:"url"`
}

type responseModel struct {
	Result string `json:"result"`
}

func (h *Handlers) APIShortenHandlerPost(w http.ResponseWriter, r *http.Request) {

	// if !h.isValidContentType(r, appJSON) {
	// 	h.publishBadRequest(&w, fmt.Errorf("content type is not as expected"))
	// 	return
	// }

	var (
		requestJSON requestModel
		i           interpreter.InterpreterJSON
	)
	if err := i.Decode(r.Body, &requestJSON); err != nil {
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

	node, err := h.service.Storage.FindByLongURL(requestJSON.URL)
	if err != nil {
		node, err = h.service.NewPairURL(requestJSON.URL)
		if err != nil {
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
