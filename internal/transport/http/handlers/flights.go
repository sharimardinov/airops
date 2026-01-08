package handlers

import (
	"airops/internal/transport/http/dto"
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

	out := make([]dto.FlightResponse, 0, len(items))
	for _, f := range items {
		out = append(out, dto.FlightFromModel(f))
	}

	writeJSON(w, http.StatusOK, out)
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

	resp := dto.FlightFromModel(item.Flight)
	writeJSON(w, http.StatusOK, resp)
}
