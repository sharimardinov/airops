// flight.go
package models

import "time"

type Flight struct {
	ID                 int64      `json:"id"`
	RouteNo            string     `json:"route_no"`
	Status             string     `json:"status"`
	DepartureAirport   string     `json:"departure_airport"`
	ArrivalAirport     string     `json:"arrival_airport"`
	ScheduledDeparture time.Time  `json:"scheduled_departure"`
	ScheduledArrival   time.Time  `json:"scheduled_arrival"`
	ActualDeparture    *time.Time `json:"actual_departure,omitempty"`
	ActualArrival      *time.Time `json:"actual_arrival,omitempty"`
}
