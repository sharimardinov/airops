package dto

import "airops/internal/domain/models"

type RouteStatResponse struct {
	From         string `json:"from"`
	To           string `json:"to"`
	FlightsCount int64  `json:"flights_count"`
}

func RouteStatFromModel(s models.RouteStat) RouteStatResponse {
	return RouteStatResponse{
		From:         s.DepartureAirport,
		To:           s.ArrivalAirport,
		FlightsCount: s.FlightsCount,
	}
}
