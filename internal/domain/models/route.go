package models

import "time"

type Route struct {
	RouteNo          string        `json:"route_no"`
	DepartureAirport string        `json:"departure_airport"`
	ArrivalAirport   string        `json:"arrival_airport"`
	AirplaneCode     string        `json:"airplane_code"`
	DaysOfWeek       []int         `json:"days_of_week"`   // [1,2,3,4,5] = пн-пт
	ScheduledTime    time.Time     `json:"scheduled_time"` // время вылета
	Duration         time.Duration `json:"duration"`       // длительность полета
}
