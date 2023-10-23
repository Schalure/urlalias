package handlers

import (
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

	
}