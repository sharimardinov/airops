// internal/domain/http/handlers/airplanes.go
package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ListAirplanes возвращает список всех самолетов
func (h *Handler) ListAirplanes(w http.ResponseWriter, r *http.Request) {
	airplanes, err := h.airplanesService.List(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, airplanes)
}

// GetAirplane получает самолет по коду
func (h *Handler) GetAirplane(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeError(w, r, fmt.Errorf("airplane code is required"))
		return
	}

	airplane, err := h.airplanesService.GetByCode(r.Context(), code)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, airplane)
}

// GetAirplaneWithSeats получает самолет с раскладкой мест
func (h *Handler) GetAirplaneWithSeats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeError(w, r, fmt.Errorf("airplane code is required"))
		return
	}

	airplane, err := h.airplanesService.GetWithSeats(r.Context(), code)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, airplane)
}

// GetAirplaneStats получает статистику по самолету
func (h *Handler) GetAirplaneStats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		writeError(w, r, fmt.Errorf("airplane code is required"))
		return
	}

	stats, err := h.airplanesService.GetStats(r.Context(), code)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
