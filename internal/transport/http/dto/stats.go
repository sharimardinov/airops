package dto

type RouteStatResponse struct {
	From         string `json:"from"`
	To           string `json:"to"`
	FlightsCount int64  `json:"flights_count"`
}
