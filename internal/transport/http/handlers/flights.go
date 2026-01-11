package handlers

import (
	"airops/internal/transport/http/dto"
	"net/http"
	"time"
)

// ListFlights возвращает список рейсов
func (h *Handler) ListFlights(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, r, err)
		return
	}

	limit := qInt(r, "limit", 100)
	offset := qInt(r, "offset", 0)

	flights, err := h.flightsService.List(r.Context(), date, limit, offset) // <-- добавь offset
	if err != nil {
		writeError(w, r, err)
		return
	}

	response := make([]dto.FlightResponse, 0, len(flights))
	for _, flight := range flights {
		response = append(response, dto.FlightFromModel(flight))
	}
	writeJSON(w, http.StatusOK, response)
}

// GetFlightByID возвращает детальную информацию о рейсе
func (h *Handler) GetFlightByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, r, err)
		return
	}

	// Получаем FlightDetails (с пассажирами)
	flightDetails, err := h.flightsService.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}

	// ✅ Конвертируем FlightDetails в DTO
	response := dto.FlightDetailsFromModel(flightDetails)

	writeJSON(w, http.StatusOK, response)
}
