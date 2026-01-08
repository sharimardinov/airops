package dto

import (
	"airops/internal/domain/models"
	"time"
)

type FlightResponse struct {
	ID                 int64      `json:"id"`
	Status             string     `json:"status"`
	DepartureAirport   string     `json:"departure_airport"`
	ArrivalAirport     string     `json:"arrival_airport"`
	ScheduledDeparture time.Time  `json:"scheduled_departure"`
	ScheduledArrival   time.Time  `json:"scheduled_arrival"`
	ActualDeparture    *time.Time `json:"actual_departure,omitempty"`
	ActualArrival      *time.Time `json:"actual_arrival,omitempty"`
}

func FlightFromModel(f models.Flight) FlightResponse {
	return FlightResponse{
		ID:                 f.ID,
		Status:             f.Status,
		DepartureAirport:   f.DepartureAirport,
		ArrivalAirport:     f.ArrivalAirport,
		ScheduledDeparture: f.ScheduledDeparture,
		ScheduledArrival:   f.ScheduledArrival,
		ActualDeparture:    f.ActualDeparture,
		ActualArrival:      f.ActualArrival,
	}
}
