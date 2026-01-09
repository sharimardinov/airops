package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Booking struct {
	BookRef     string
	BookDate    time.Time
	TotalAmount decimal.Decimal
	Tickets     []Ticket
}
