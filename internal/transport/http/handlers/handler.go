package handlers

import "airops/internal/usecase"

type Handler struct {
	flights    *usecase.FlightsService
	passengers *usecase.PassengersService
	stats      *usecase.StatsRoutesService
	health     *usecase.HealthService
}

func New(
	flights *usecase.FlightsService,
	passengers *usecase.PassengersService,
	stats *usecase.StatsRoutesService,
	health *usecase.HealthService,
) *Handler {
	return &Handler{
		flights:    flights,
		passengers: passengers,
		stats:      stats,
		health:     health,
	}
}
