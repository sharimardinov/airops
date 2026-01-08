package handlers

import (
	"net/http"
	"time"
)

func (h *Handler) TopRoutes(w http.ResponseWriter, r *http.Request) {
	limit := qInt(r, "limit", 10)

	var from, to time.Time

	items, err := h.stats.TopRoutes(r.Context(), from, to, limit)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, items)
}
