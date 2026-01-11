package handlers

import (
// если хочешь типы в ответах — можешь импортить:
// "airops/internal/transport/http/dto"
// "airops/internal/domain/models"
)

// ---- GENERAL (один раз на проект) ----

// @title airops API
// @version 1.0
// @description Demo API for flights/bookings/stats
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func _swaggerGeneral() {}

// ---- AIRPLANES ----

// ListAirplanes godoc
// @Summary List airplanes
// @Tags airplanes
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /airplanes [get]
func _swaggerListAirplanes() {}

// GetAirplane godoc
// @Summary Get airplane by code
// @Tags airplanes
// @Produce json
// @Param code path string true "Airplane code"
// @Security ApiKeyAuth
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /airplanes/{code} [get]
func _swaggerGetAirplane() {}

// GetAirplaneWithSeats godoc
// @Summary Get airplane with seats
// @Tags airplanes
// @Produce json
// @Param code path string true "Airplane code"
// @Security ApiKeyAuth
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /airplanes/{code}/seats [get]
func _swaggerGetAirplaneWithSeats() {}

// GetAirplaneStats godoc
// @Summary Get airplane stats
// @Tags airplanes
// @Produce json
// @Param code path string true "Airplane code"
// @Security ApiKeyAuth
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /airplanes/{code}/stats [get]
func _swaggerGetAirplaneStats() {}

// ---- AIRPORTS ----

// ListAirports godoc
// @Summary List airports
// @Tags airports
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /airports [get]
func _swaggerListAirports() {}

// SearchAirports godoc
// @Summary Search airports
// @Tags airports
// @Produce json
// @Param city query string false "City"
// @Param country query string false "Country"
// @Security ApiKeyAuth
// @Success 200 {array} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /airports/search [get]
func _swaggerSearchAirports() {}

// GetAirport godoc
// @Summary Get airport by code
// @Tags airports
// @Produce json
// @Param code path string true "Airport code"
// @Security ApiKeyAuth
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /airports/{code} [get]
func _swaggerGetAirport() {}

// ---- FLIGHTS ----

// ListFlights godoc
// @Summary List flights
// @Tags flights
// @Produce json
// @Param from query string false "From date YYYY-MM-DD"
// @Param to query string false "To date YYYY-MM-DD"
// @Param limit query int false "Limit"
// @Security ApiKeyAuth
// @Success 200 {array} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /flights [get]
func _swaggerListFlights() {}

// GetFlightByID godoc
// @Summary Get flight by id
// @Tags flights
// @Produce json
// @Param id path int true "Flight ID"
// @Security ApiKeyAuth
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /flights/{id} [get]
func _swaggerGetFlightByID() {}

// SearchFlights godoc
// @Summary Search flights
// @Tags flights
// @Produce json
// @Param from query string true "From date YYYY-MM-DD"
// @Param to query string true "To date YYYY-MM-DD"
// @Param from_airport query string false "From airport code"
// @Param to_airport query string false "To airport code"
// @Param limit query int false "Limit"
// @Security ApiKeyAuth
// @Success 200 {array} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /flights/search [get]
func _swaggerSearchFlights() {}

// ListFlightPassengers godoc
// @Summary List passengers of flight
// @Tags flights
// @Produce json
// @Param id path int true "Flight ID"
// @Param limit query int false "Limit"
// @Security ApiKeyAuth
// @Success 200 {array} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /flights/{id}/passengers [get]
func _swaggerListFlightPassengers() {}

// ---- BOOKINGS ----

// CreateBooking godoc
// @Summary Create booking
// @Tags bookings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body object true "Booking request"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 409 {object} object
// @Failure 500 {object} object
// @Router /bookings [post]
func _swaggerCreateBooking() {}

// GetBooking godoc
// @Summary Get booking
// @Tags bookings
// @Produce json
// @Param bookRef path string true "Booking reference"
// @Security ApiKeyAuth
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /bookings/{bookRef} [get]
func _swaggerGetBooking() {}

// CancelBooking godoc
// @Summary Cancel booking
// @Tags bookings
// @Produce json
// @Param bookRef path string true "Booking reference"
// @Security ApiKeyAuth
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /bookings/{bookRef} [delete]
func _swaggerCancelBooking() {}

// ---- PASSENGERS ----

// GetPassengerBookings godoc
// @Summary Get passenger bookings
// @Tags passengers
// @Produce json
// @Param passengerID path string true "Passenger ID"
// @Security ApiKeyAuth
// @Success 200 {array} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /passengers/{passengerID}/bookings [get]
func _swaggerGetPassengerBookings() {}

// ---- STATS ----

// TopRoutes godoc
// @Summary Top routes stats
// @Tags stats
// @Produce json
// @Param from query string true "From date YYYY-MM-DD"
// @Param to query string true "To date YYYY-MM-DD"
// @Param limit query int false "Limit"
// @Security ApiKeyAuth
// @Success 200 {array} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /stats/routes [get]
func _swaggerTopRoutes() {}
