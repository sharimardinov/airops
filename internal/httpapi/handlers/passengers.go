package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListFlightPassengers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad flight id"})
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

	items, err := h.passengers.ListByFlight(r.Context(), id, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"flight_id": id,
		"items":     items,
		"limit":     limit,
		"offset":    offset,
	})
}
