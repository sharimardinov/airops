package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type BookingDetails struct {
	BookRef     string          `json:"book_ref"`
	BookDate    time.Time       `json:"book_date"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Tickets     []TicketDetails `json:"tickets"`
}

type TicketDetails struct {
	TicketNo      string       `json:"ticket_no"`
	PassengerID   string       `json:"passenger_id"`
	PassengerName string       `json:"passenger_name"`
	Flights       []FlightInfo `json:"flights"`
}

type FlightInfo struct {
	FlightID           int64     `json:"flight_id"`
	RouteNo            string    `json:"route_no"`
	Status             string    `json:"status"`
	FareConditions     string    `json:"fare_conditions"`
	SeatNo             string    `json:"seat_no,omitempty"`
	BoardingNo         int       `json:"boarding_no,omitempty"`
	ScheduledDeparture time.Time `json:"scheduled_departure"`
	ScheduledArrival   time.Time `json:"scheduled_arrival"`
}
