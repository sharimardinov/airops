package handlers

import (
	"airops/internal/transport/http/dto"
	"net/http"
	"time"
)

func (h *Handler) TopRoutes(w http.ResponseWriter, r *http.Request) {
	limit := qInt(r, "limit", 10)

	fromStr := r.URL.Query().Get("from") // YYYY-MM-DD
	toStr := r.URL.Query().Get("to")     // YYYY-MM-DD

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{Error: "bad from (use YYYY-MM-DD)"})
			return
		}
	}
	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{Error: "bad to (use YYYY-MM-DD)"})
			return
		}
		// включим весь день "to"
		to = to.Add(24 * time.Hour)
	}

	items, err := h.statsService.TopRoutes(r.Context(), from, to, limit)
	if err != nil {
		writeError(w, r, err)
		return
	}

	out := make([]dto.RouteStatResponse, 0, len(items))
	for _, s := range items {
		out = append(out, dto.RouteStatFromModel(s))
	}

	writeJSON(w, http.StatusOK, out)
}
