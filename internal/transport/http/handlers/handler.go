package handlers

import "airops/internal/usecase"

type Handler struct {
	flights    *usecase.FlightsService
	passengers *usecase.PassengersService
	stats      *usecase.StatsRoutesService
}

func New(
	flights *usecase.FlightsService,
	passengers *usecase.PassengersService,
	stats *usecase.StatsRoutesService,
) *Handler {
	return &Handler{
		flights:    flights,
		passengers: passengers,
		stats:      stats,
	}
}
