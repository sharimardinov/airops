package dto

import "airops/internal/domain/models"

type CreateBookingRequest struct {
	FlightID      int      `json:"flight_id" validate:"required"`
	PassengerName string   `json:"passenger_name" validate:"required"`
	PassengerID   string   `json:"passenger_id" validate:"required"`
	Seats         []string `json:"seats" validate:"required,min=1"`
	FareClass     string   `json:"fare_class" validate:"required,oneof=Economy Comfort Business"`
}

type BookingResponse struct {
	BookRef     string          `json:"book_ref"`
	BookDate    string          `json:"book_date"`
	TotalAmount float64         `json:"total_amount"`
	Status      string          `json:"status"`
	Tickets     []TicketDetails `json:"tickets"`
}

type TicketDetails struct {
	TicketNo      string `json:"ticket_no"`
	PassengerName string `json:"passenger_name"`
	SeatNo        string `json:"seat_no"`
	FareClass     string `json:"fare_class"`
}

func FromBooking(b *models.Booking) BookingResponse {
	return BookingResponse{
		BookRef: b.BookRef}
}
