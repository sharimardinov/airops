package handlers

import (
	"airops/internal/transport/http/dto"
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

	out := make([]dto.RouteStatResponse, 0, len(items))
	for _, s := range items {
		out = append(out, dto.RouteStatFromModel(s))
	}

	writeJSON(w, http.StatusOK, out)
}
