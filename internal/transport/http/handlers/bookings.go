package handlers

import (
	"airops/internal/domain/models"
	"airops/internal/transport/http/dto"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid json"})
		return
	}

	if err := validateStruct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "validation failed"})
		return
	}

	domainReq := models.BookingRequest{
		FlightID:      int64(req.FlightID),
		PassengerName: req.PassengerName,
		PassengerID:   req.PassengerID,
		Seats:         req.Seats,
		FareClass:     req.FareClass,
	}

	details, err := h.bookingService.Create(r.Context(), domainReq)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusCreated, details)
}

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	bookRef := chi.URLParam(r, "bookRef")
	if bookRef == "" {
		writeError(w, r, fmt.Errorf("bookRef is required"))
		return
	}

	details, err := h.bookingService.GetByRef(r.Context(), bookRef)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, details)
}

// GetPassengerBookings получает все бронирования пассажира
func (h *Handler) GetPassengerBookings(w http.ResponseWriter, r *http.Request) {
	passengerID := chi.URLParam(r, "passengerID")
	if passengerID == "" {
		writeError(w, r, fmt.Errorf("passenger ID is required"))
		return
	}

	bookings, err := h.bookingService.GetByPassenger(r.Context(), passengerID)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, bookings)
}

// CancelBooking отменяет бронирование
func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	bookRef := chi.URLParam(r, "bookRef")
	if bookRef == "" {
		writeError(w, r, fmt.Errorf("booking reference is required"))
		return
	}

	if err := h.bookingService.Cancel(r.Context(), bookRef); err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "booking cancelled successfully",
	})
}
