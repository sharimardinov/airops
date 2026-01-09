package models

import "time"

type BoardingPass struct {
	TicketNo     string
	FlightID     int
	SeatNo       string
	BoardingNo   int
	BoardingTime time.Time
}
