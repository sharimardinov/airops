package handlers

import (
	"airops/internal/store"
	"net/http"
	"time"
)

func (h *Handler) RoutesStats(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := store.RoutesStatsFilter{
		Limit:  qInt(r, "limit", 50),
		Offset: qInt(r, "offset", 0),
		Sort:   q.Get("sort"),
	}

	if from := q.Get("from"); from != "" {
		filter.From = &from
	}
	if to := q.Get("to"); to != "" {
		filter.To = &to
	}
	if dateFrom := q.Get("date_from"); dateFrom != "" {
		if t, err := time.Parse(time.RFC3339, dateFrom); err == nil {
			filter.DateFrom = &t
		}
	}
	if dateTo := q.Get("date_to"); dateTo != "" {
		if t, err := time.Parse(time.RFC3339, dateTo); err == nil {
			filter.DateTo = &t
		}
	}

	items, err := h.stats.Routes(r.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, items)
}
