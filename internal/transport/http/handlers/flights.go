package handlers

import (
	"net/http"
	"time"
)

func (h *Handler) ListFlights(w http.ResponseWriter, r *http.Request) {
	limit := qInt(r, "limit", 100)
	offset := qInt(r, "offset", 0)

	var from, to time.Time

	items, err := h.flights.List(r.Context(), from, to, limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) GetFlightByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad id"})
		return
	}

	item, err := h.flights.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, item)
}
