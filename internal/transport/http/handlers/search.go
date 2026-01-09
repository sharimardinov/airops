package handlers

import (
	"airops/internal/domain/models"
	"fmt"
	"net/http"
	"time"
)

// SearchFlights обрабатывает поиск рейсов
func (h *Handler) SearchFlights(w http.ResponseWriter, r *http.Request) {
	// Парсим параметры поиска
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	dateStr := r.URL.Query().Get("date")
	passengers := qInt(r, "passengers", 1)
	fareClass := r.URL.Query().Get("fare_class")

	if from == "" || to == "" || dateStr == "" {
		writeError(w, r, fmt.Errorf("missing required parameters: from, to, date"))
		return
	}

	// Парсим дату
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, r, fmt.Errorf("invalid date format (use YYYY-MM-DD)"))
		return
	}

	// По умолчанию Economy
	if fareClass == "" {
		fareClass = "Economy"
	}

	params := models.FlightSearchParams{
		DepartureAirport: from,
		ArrivalAirport:   to,
		DepartureDate:    date,
		Passengers:       passengers,
		FareClass:        fareClass,
	}

	// Выполняем поиск
	results, err := h.searchService.SearchFlights(r.Context(), params)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, results)
}

// internal/domain/http/handlers/airports.go
