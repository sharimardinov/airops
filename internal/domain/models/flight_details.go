// flight_details.go
package models

type FlightDetails struct {
	Flight     Flight            `json:"flight"`
	Passengers []FlightPassenger `json:"passengers"`
}
