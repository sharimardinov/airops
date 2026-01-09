package models

import "time"

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
