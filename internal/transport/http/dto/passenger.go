package dto

import (
	"airops/internal/domain/models"
	"time"
)

type FlightPassengerResponse struct {
	TicketNo     string     `json:"ticket_no"`
	Name         string     `json:"name"`
	FareClass    string     `json:"fare_class"`
	Boarded      bool       `json:"boarded"`
	SeatNo       *string    `json:"seat_no,omitempty"`
	BoardingTime *time.Time `json:"boarding_time,omitempty"`
}

func PassengerFromModel(p models.FlightPassenger) FlightPassengerResponse {
	return FlightPassengerResponse{
		TicketNo:     p.TicketNo,
		Name:         p.PassengerName,
		FareClass:    p.FareConditions,
		Boarded:      p.Boarded,
		SeatNo:       p.SeatNo,
		BoardingTime: p.BoardingTime,
	}
}
