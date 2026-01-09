package models

type RouteStat struct {
	DepartureAirport string `json:"departure_airport"`
	ArrivalAirport   string `json:"arrival_airport"`
	FlightsCount     int64  `json:"flights_count"`
}
