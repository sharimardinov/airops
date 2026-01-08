package handlers

import (
	"airops/internal/transport/http/dto"
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

	out := make([]dto.FlightPassengerResponse, 0, len(items))
	for _, p := range items {
		out = append(out, dto.PassengerFromModel(p))
	}
	writeJSON(w, http.StatusOK, out)
}
