// flight_details.go
package models

import "time"

type FlightDetails struct {
	FlightID           int64     `json:"flight_id"`
	RouteNo            string    `json:"route_no"`
	Status             string    `json:"status"`
	ScheduledDeparture time.Time `json:"scheduled_departure"`
	ScheduledArrival   time.Time `json:"scheduled_arrival"`
	ActualDeparture    time.Time `json:"actual_departure,omitempty"`
	ActualArrival      time.Time `json:"actual_arrival,omitempty"`

	// Информация о маршруте
	DepartureAirport     string `json:"departure_airport"`
	ArrivalAirport       string `json:"arrival_airport"`
	DepartureAirportName string `json:"departure_airport_name"`
	DepartureCity        string `json:"departure_city"`
	ArrivalAirportName   string `json:"arrival_airport_name"`
	ArrivalCity          string `json:"arrival_city"`

	// Информация о самолете
	AirplaneCode  string        `json:"airplane_code"`
	AirplaneModel string        `json:"airplane_model,omitempty"`
	Duration      time.Duration `json:"duration,omitempty"`

	// Пассажиры (опционально)
	Passengers []FlightPassenger `json:"passengers,omitempty"`
}
