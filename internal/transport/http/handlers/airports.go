package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ListAirports возвращает список всех аэропортов
func (h *Handler) ListAirports(w http.ResponseWriter, r *http.Request) {
	airports, err := h.airportsService.List(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, airports)
}

// SearchAirports ищет аэропорты по названию города
func (h *Handler) SearchAirports(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		writeError(w, r, fmt.Errorf("city parameter is required"))
		return
	}

	airports, err := h.airportsService.SearchByCity(r.Context(), city)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, airports)
}

// GetAirport получает аэропорт по коду
func (h *Handler) GetAirport(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeError(w, r, fmt.Errorf("airport code is required"))
		return
	}

	airport, err := h.airportsService.GetByCode(r.Context(), code)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, airport)
}
