package handlers

import (
	usecase2 "airops/internal/app/usecase"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool              *pgxpool.Pool
	flightsService    *usecase2.FlightsService
	passengersService *usecase2.PassengersService
	statsService      *usecase2.StatsRoutesService
	healthService     *usecase2.HealthService
	bookingService    *usecase2.BookingService
	searchService     *usecase2.SearchService
	airportsService   *usecase2.AirportsService
	airplanesService  *usecase2.AirplanesService
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func New(
	pool *pgxpool.Pool,
	flightsService *usecase2.FlightsService,
	passengersService *usecase2.PassengersService,
	statsService *usecase2.StatsRoutesService,
	healthService *usecase2.HealthService,
	bookingService *usecase2.BookingService,
	searchService *usecase2.SearchService,
	airportsService *usecase2.AirportsService,
	airplanesService *usecase2.AirplanesService, // ✨ NEW
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
