package handlers

import (
	"airops/internal/service"
	"airops/internal/store"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	flights    *service.FlightsService
	passengers *service.PassengersService
	stats      *service.StatsRoutesService
}

func New(pool *pgxpool.Pool) *Handler {
	return &Handler{
		flights: service.NewFlightsService(
			store.NewFlightsStore(pool),
			store.NewFlightDetailsStore(pool),
		),
		passengers: service.NewPassengersService(
			store.NewPassengersStore(pool),
		),
		stats: service.NewStatsRoutesService(
			store.NewStatsStore(pool),
		),
	}
}
