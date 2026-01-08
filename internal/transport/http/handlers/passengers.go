package handlers

import (
	"net/http"
)

func (h *Handler) ListFlightPassengers(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad id"})
		return
	}

	limit := qInt(r, "limit", 100)
	offset := qInt(r, "offset", 0)

	items, err := h.passengers.ListByFlightID(r.Context(), id, limit, offset)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, items)
}
