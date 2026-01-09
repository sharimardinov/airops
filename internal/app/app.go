// internal/app/app.go
package app

import (
	usecase2 "airops/internal/application/usecase"
	repositories2 "airops/internal/infrastructure/postgres/repositories"
	"context"
	"fmt"
	"net/http"
	"time"

	transporthttp "airops/internal/transport/http"
	"airops/internal/transport/http/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	server *http.Server
}

func New(pool *pgxpool.Pool, addr string) *App {
	// Repositories
	flightsRepo := repositories2.NewFlightsRepo(pool)
	passengersRepo := repositories2.NewPassengersRepo(pool)
	statsRoutesRepo := repositories2.NewStatsRoutesRepo(pool)
	healthRepo := repositories2.NewHealthRepo(pool)
	bookingsRepo := repositories2.NewBookingsRepo(pool)
	seatsRepo := repositories2.NewSeatsRepo(pool)
	ticketsRepo := repositories2.NewTicketsRepo(pool)
	airportsRepo := repositories2.NewAirportsRepo(pool)

	// Usecases
	flightsService := usecase2.NewFlightsService(flightsRepo, passengersRepo)
	passengersService := usecase2.NewPassengersService(passengersRepo)
	statsService := usecase2.NewStatsRoutesService(statsRoutesRepo)
	healthService := usecase2.NewHealthService(healthRepo)
	bookingService := usecase2.NewBookingService(bookingsRepo, flightsRepo, seatsRepo, ticketsRepo)
	searchService := usecase2.NewSearchService(flightsRepo, seatsRepo)
	airportsService := usecase2.NewAirportsService(airportsRepo)

	// Handlers
	h := handlers.New(
		flightsService,
		passengersService,
		statsService,
		healthService,
		bookingService,
		searchService,
		airportsService,
	)

	// Router
	router := transporthttp.New(h)

	// Server
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{server: server}
}

func (a *App) Run() error {
	fmt.Printf("Starting server on %s\n", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
