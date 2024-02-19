package handlers

import (
	"net/http"
)

func (h *Handler) PingGet(w http.ResponseWriter, r *http.Request) {

	if !h.service.Storage.IsConnected() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
