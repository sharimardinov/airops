// internal/app/app.go
package app

import (
	usecase2 "airops/internal/app/usecase"
	"airops/internal/infrastructure/postgres/repositories"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	transporthttp "airops/internal/transport/http"
	"airops/internal/transport/http/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
	//"golang.org/x/tools/go/cfg"
)

type App struct {
	server *http.Server
	pool   *pgxpool.Pool
}

func New(pool *pgxpool.Pool, addr string) *App {
	flightsRepo := repositories.NewFlightsRepo(pool)
	passengersRepo := repositories.NewPassengersRepo(pool)
	statsRoutesRepo := repositories.NewStatsRoutesRepo(pool)
	healthRepo := repositories.NewHealthRepo(pool)
	bookingsRepo := repositories.NewBookingsRepo(pool)
	seatsRepo := repositories.NewSeatsRepo(pool)
	ticketsRepo := repositories.NewTicketsRepo(pool)
	airportsRepo := repositories.NewAirportsRepo(pool)
	airplanesRepo := repositories.NewAirplanesRepo(pool) // ✨ NEW

	flightsService := usecase2.NewFlightsService(flightsRepo, passengersRepo)
	passengersService := usecase2.NewPassengersService(passengersRepo)
	statsService := usecase2.NewStatsRoutesService(statsRoutesRepo)
	healthService := usecase2.NewHealthService(healthRepo)
	bookingService := usecase2.NewBookingService(bookingsRepo, flightsRepo, seatsRepo, ticketsRepo)
	searchService := usecase2.NewSearchService(flightsRepo, seatsRepo)
	airportsService := usecase2.NewAirportsService(airportsRepo)
	airplanesService := usecase2.NewAirplanesService(airplanesRepo) // ✨ NEW

	h := handlers.New(
		pool,
		flightsService,
		passengersService,
		statsService,
		healthService,
		bookingService,
		searchService,
		airportsService,
		airplanesService,
	)

	router := transporthttp.New(h)

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 2 * time.Minute,
		IdleTimeout:  2 * time.Minute,
	}

	return &App{
		server: server,
		pool:   pool,
	}
}

func (a *App) Run() error {
	go func() {
		log.Printf("Starting server on %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}
	a.pool.Close()

	return nil
}
