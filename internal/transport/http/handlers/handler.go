// internal/domain/http/handlers/handler.go
package handlers

import (
	"airops/internal/usecase"
)

type Handler struct {
	flightsService    *usecase.FlightsService
	passengersService *usecase.PassengersService
	statsService      *usecase.StatsRoutesService
	healthService     *usecase.HealthService
	bookingService    *usecase.BookingService
	searchService     *usecase.SearchService
	airportsService   *usecase.AirportsService
}

func New(
	flightsService *usecase.FlightsService,
	passengersService *usecase.PassengersService,
	statsService *usecase.StatsRoutesService,
	healthService *usecase.HealthService,
	bookingService *usecase.BookingService,
	searchService *usecase.SearchService,
	airportsService *usecase.AirportsService,
) *Handler {
	return &Handler{
		flightsService:    flightsService,
		passengersService: passengersService,
		statsService:      statsService,
		healthService:     healthService,
		bookingService:    bookingService,
		searchService:     searchService,
		airportsService:   airportsService,
	}
}
