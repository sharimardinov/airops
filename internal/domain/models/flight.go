// flight.go
package models

import "time"

type Flight struct {
	FlightID           int64     `json:"flight_id"`
	RouteNo            string    `json:"route_no"`
	Status             string    `json:"status"`
	ScheduledDeparture time.Time `json:"scheduled_departure"`
	ScheduledArrival   time.Time `json:"scheduled_arrival"`
	ActualDeparture    time.Time `json:"actual_departure,omitempty"`
	ActualArrival      time.Time `json:"actual_arrival,omitempty"`
}
