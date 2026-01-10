package models

type FlightSearchResult struct {
	FlightDetails          // Встраиваем все поля FlightDetails
	AvailableSeats int     `json:"available_seats"`
	Price          float64 `json:"price"`
}
