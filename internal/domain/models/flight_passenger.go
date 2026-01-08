// flight_passenger.go
package models

import "time"

// FlightPassenger — “пассажир на рейсе” (не абстрактный пассажир как человек).
type FlightPassenger struct {
	TicketNo       string     `json:"ticket_no"`
	PassengerName  string     `json:"passenger_name"`
	FareConditions string     `json:"fare_conditions"`
	Boarded        bool       `json:"boarded"`
	SeatNo         *string    `json:"seat_no,omitempty"`
	BoardingTime   *time.Time `json:"boarding_time,omitempty"`
}
