package models

import "github.com/shopspring/decimal"

type TicketSegment struct {
	TicketNo       string          `json:"ticket_no"`
	FlightID       int64           `json:"flight_id"`
	FareConditions string          `json:"fare_conditions"` // Economy, Comfort, Business
	Price          decimal.Decimal `json:"price"`
}
