package handlers

import (
	"airops/internal/store"
	"context"
	"net/http"
	"time"
)

type FlightsRepo interface {
	List(ctx context.Context, from, to time.Time, limit, offset int) ([]store.Flight, error)
}

type FlightDetailsRepo interface {
	GetByID(ctx context.Context, id int64) (store.FlightDetails, error)
}

type FlightsService struct {
	flights FlightsRepo
	details FlightDetailsRepo
}

func NewFlightsService(f FlightsRepo, d FlightDetailsRepo) *FlightsService {
	return &FlightsService{flights: f, details: d}
}

func (h *Handler) ListFlights(w http.ResponseWriter, r *http.Request) {
	limit := qInt(r, "limit", 100)
	offset := qInt(r, "offset", 0)

	// Исправлено: добавлены параметры from и to
	items, err := h.flights.List(r.Context(), time.Time{}, time.Time{}, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, items)
}
