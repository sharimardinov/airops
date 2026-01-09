// internal/domain/models/booking_details.go
package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// BookingDetails - полная информация о бронировании с билетами и рейсами
type BookingDetails struct {
	BookRef     string          `json:"book_ref"`
	BookDate    time.Time       `json:"book_date"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Tickets     []TicketDetails `json:"tickets"`
}

// TicketDetails - детальная информация о билете
type TicketDetails struct {
	TicketNo      string       `json:"ticket_no"`
	PassengerID   string       `json:"passenger_id"`
	PassengerName string       `json:"passenger_name"`
	Flights       []FlightInfo `json:"flights"`
}

// FlightInfo - информация о рейсе в билете
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
