package models

type AirplaneWithSeats struct {
	Airplane Airplane `json:"airplane"`
	Seats    []Seat   `json:"seats"`
}
