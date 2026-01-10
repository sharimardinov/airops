package models

import (
	"errors"
	"time"
)

type FlightSearchParams struct {
	DepartureAirport string    `json:"departure_airport"`
	ArrivalAirport   string    `json:"arrival_airport"`
	DepartureDate    time.Time `json:"departure_date"`
	Passengers       int       `json:"passengers"`
	FareClass        string    `json:"fare_class"` // Economy, Comfort, Business
}

// BookingRequest - запрос на создание бронирования
type BookingRequest struct {
	FlightID      int64    `json:"flight_id"`
	PassengerName string   `json:"passenger_name"`
	PassengerID   string   `json:"passenger_id"`
	Seats         []string `json:"seats"`      // ["1A", "1B"]
	FareClass     string   `json:"fare_class"` // Economy, Comfort, Business
}

func (p *FlightSearchParams) Validate() error {
	if p.DepartureAirport == "" {
		return errors.New("departure_airport is required")
	}
	if p.ArrivalAirport == "" {
		return errors.New("arrival_airport is required")
	}
	if p.DepartureDate.IsZero() {
		return errors.New("departure_date is required")
	}
	if p.Passengers <= 0 {
		return errors.New("passengers must be positive")
	}
	return nil
}
