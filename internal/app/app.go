// internal/app/app.go
package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"airops/internal/infra/db/pg/repo"
	transporthttp "airops/internal/transport/http"
	"airops/internal/transport/http/handlers"
	"airops/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	server *http.Server
}

func New(pool *pgxpool.Pool, addr string) *App {
	// Repositories
	flightsRepo := repo.NewFlightsRepo(pool)
	passengersRepo := repo.NewPassengersRepo(pool)
	statsRoutesRepo := repo.NewStatsRoutesRepo(pool)
	healthRepo := repo.NewHealthRepo(pool)
	bookingsRepo := repo.NewBookingsRepo(pool)
	seatsRepo := repo.NewSeatsRepo(pool)
	ticketsRepo := repo.NewTicketsRepo(pool)
	airportsRepo := repo.NewAirportsRepo(pool)

	// Usecases
	flightsService := usecase.NewFlightsService(flightsRepo, passengersRepo)
	passengersService := usecase.NewPassengersService(passengersRepo)
	statsService := usecase.NewStatsRoutesService(statsRoutesRepo)
	healthService := usecase.NewHealthService(healthRepo)
	bookingService := usecase.NewBookingService(bookingsRepo, flightsRepo, seatsRepo, ticketsRepo)
	searchService := usecase.NewSearchService(flightsRepo, seatsRepo)
	airportsService := usecase.NewAirportsService(airportsRepo)

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
