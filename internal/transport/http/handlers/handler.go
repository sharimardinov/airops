// internal/domain/http/handlers/handler.go
package handlers

import (
	usecase2 "airops/internal/application/usecase"
)

type Handler struct {
	flightsService    *usecase2.FlightsService
	passengersService *usecase2.PassengersService
	statsService      *usecase2.StatsRoutesService
	healthService     *usecase2.HealthService
	bookingService    *usecase2.BookingService
	searchService     *usecase2.SearchService
	airportsService   *usecase2.AirportsService
}

func New(
	flightsService *usecase2.FlightsService,
	passengersService *usecase2.PassengersService,
	statsService *usecase2.StatsRoutesService,
	healthService *usecase2.HealthService,
	bookingService *usecase2.BookingService,
	searchService *usecase2.SearchService,
	airportsService *usecase2.AirportsService,
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
