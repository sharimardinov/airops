package handlers

import (
	"net/http"
	"time"

	"airops/internal/store"
)

func (h *Handler) RoutesStats(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	var fromPtr *string
	if from != "" {
		fromPtr = &from
	}
	var toPtr *string
	if to != "" {
		toPtr = &to
	}

	dateFrom, err := parseDateParam(r, "date_from")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad date_from (use YYYY-MM-DD)"})
		return
	}
	dateTo, err := parseDateParam(r, "date_to")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad date_to (use YYYY-MM-DD)"})
		return
	}

	// дефолт: последние 30 дней
	if dateFrom == nil && dateTo == nil {
		now := time.Now().UTC()
		df := now.AddDate(0, 0, -30)
		dt := now
		dateFrom = &df
		dateTo = &dt
	}

	// защита от идиотских диапазонов
	if dateFrom != nil && dateTo != nil && !dateFrom.Before(*dateTo) {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "date_from must be < date_to"})
		return
	}

	limit, err := parseIntParam(r, "limit", 50, 1, 200)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad limit"})
		return
	}
	offset, err := parseIntParam(r, "offset", 0, 0, 1_000_000_000)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad offset"})
		return
	}

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "boarded"
	}
	if sort != "boarded" && sort != "load" {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid sort, allowed: boarded, load"})
		return
	}

	items, err := h.stats.Routes(r.Context(), store.RoutesStatsFilter{
		From:     fromPtr,
		To:       toPtr,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Limit:    limit,
		Offset:   offset,
		Sort:     sort,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":     items,
		"date_from": dateFrom,
		"date_to":   dateTo,
		"limit":     limit,
		"offset":    offset,
		"sort":      sort,
	})
}
