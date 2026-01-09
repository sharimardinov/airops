package models

// FlightSearchResult - результат поиска рейса с дополнительной информацией
type FlightSearchResult struct {
	Flight         Flight  `json:"flight"`
	AvailableSeats int     `json:"available_seats"`
	Price          float64 `json:"price"`
}
