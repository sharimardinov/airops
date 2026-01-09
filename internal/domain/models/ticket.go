package models

type Ticket struct {
	TicketNo      string
	BookRef       string
	PassengerID   string
	PassengerName string
	Outbound      bool
}
