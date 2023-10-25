package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type request struct{
	URL string `json:"url"`
}

type response struct{
	Result string `json:"result"`
}

func (h *Handlers) ApiShortenHandlerPost(w http.ResponseWriter, r *http.Request){
	
	if !h.isValidContentType(r, appJSON){
		h.publishBadRequest(&w, fmt.Errorf("content type is not as expected"))
		return
	}

	var requestJson request
	if err := json.NewDecoder(r.Body).Decode(&requestJson); err != nil {
		h.publishBadRequest(&w, fmt.Errorf("can't decode JSON content"))
		data, _ := io.ReadAll(r.Body)
		h.logger.Infow(
			"Can't decode JSON content",
			"Content", data,
		)
		return
	}

	if !h.isValidURL(requestJson.URL){
		h.publishBadRequest(&w, fmt.Errorf("url is not in the correct format"))
		return
	}

	node, err := h.service.Storage.FindByLongURL(requestJson.URL)
	if err != nil {
		node, err = h.service.NewPairURL(requestJson.URL); if err != nil{
			h.publishBadRequest(&w, err)
			return
		}
	}

	var buf bytes.Buffer
	var resp = response{
		Result: h.config.BaseURL() + "/" + node.ShortKey,
	}
	if err := json.NewEncoder(&buf).Encode(&resp); err != nil{
		h.publishBadRequest(&w, err)
		return
	}
	h.logger.Infow(
		"Serch/Create alias key",
		"Long URL", node.LongURL,
		"Alias URL", resp.Result,
	)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}