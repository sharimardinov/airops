package handlers

import (
	"airops/internal/application/usecase"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool              *pgxpool.Pool
	flightsService    *usecase.FlightsService
	passengersService *usecase.PassengersService
	statsService      *usecase.StatsRoutesService
	healthService     *usecase.HealthService
	bookingService    *usecase.BookingService
	searchService     *usecase.SearchService
	airportsService   *usecase.AirportsService
	airplanesService  *usecase.AirplanesService
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func New(
	pool *pgxpool.Pool,
	flightsService *usecase.FlightsService,
	passengersService *usecase.PassengersService,
	statsService *usecase.StatsRoutesService,
	healthService *usecase.HealthService,
	bookingService *usecase.BookingService,
	searchService *usecase.SearchService,
	airportsService *usecase.AirportsService,
	airplanesService *usecase.AirplanesService, // ✨ NEW
) *Handler {
	return &Handler{
		pool:              pool,
		flightsService:    flightsService,
		passengersService: passengersService,
		statsService:      statsService,
		healthService:     healthService,
		bookingService:    bookingService,
		searchService:     searchService,
		airportsService:   airportsService,
		airplanesService:  airplanesService, // ✨ NEW
	}
}
