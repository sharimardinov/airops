package dto

import (
	"airops/internal/domain/models"
	"time"
)

type FlightResponse struct {
	FlightID           int64      `json:"flight_id"` // ✅ было ID
	RouteNo            string     `json:"route_no"`
	Status             string     `json:"status"`
	DepartureAirport   string     `json:"departure_airport"` // ✅ добавлено
	ArrivalAirport     string     `json:"arrival_airport"`   // ✅ добавлено
	ScheduledDeparture time.Time  `json:"scheduled_departure"`
	ScheduledArrival   time.Time  `json:"scheduled_arrival"`
	ActualDeparture    *time.Time `json:"actual_departure,omitempty"`
	ActualArrival      *time.Time `json:"actual_arrival,omitempty"`
}

// FlightFromModel конвертирует Flight в DTO
func FlightFromModel(f models.Flight) FlightResponse {
	resp := FlightResponse{
		FlightID:           f.FlightID, // ✅ было f.ID
		RouteNo:            f.RouteNo,
		Status:             f.Status,
		ScheduledDeparture: f.ScheduledDeparture,
		ScheduledArrival:   f.ScheduledArrival,
		// DepartureAirport и ArrivalAirport отсутствуют в базовом Flight
		// Они есть только в FlightDetails
	}

	// ✅ Правильная обработка указателей
	if !f.ActualDeparture.IsZero() {
		resp.ActualDeparture = &f.ActualDeparture
	}
	if !f.ActualArrival.IsZero() {
		resp.ActualArrival = &f.ActualArrival
	}

	return resp
}

// FlightDetailsFromModel конвертирует FlightDetails в DTO
func FlightDetailsFromModel(f models.FlightDetails) FlightResponse {
	resp := FlightResponse{
		FlightID:           f.FlightID,
		RouteNo:            f.RouteNo,
		Status:             f.Status,
		DepartureAirport:   f.DepartureAirport, // ✅ теперь есть
		ArrivalAirport:     f.ArrivalAirport,   // ✅ теперь есть
		ScheduledDeparture: f.ScheduledDeparture,
		ScheduledArrival:   f.ScheduledArrival,
	}

	if !f.ActualDeparture.IsZero() {
		resp.ActualDeparture = &f.ActualDeparture
	}
	if !f.ActualArrival.IsZero() {
		resp.ActualArrival = &f.ActualArrival
	}

	return resp
}
