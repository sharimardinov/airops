package handlers

import (
	"net/http"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.health == nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "health service not wired"})
		return
	}

	if err := h.health.Ready(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, apiError{Error: "not ready"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
