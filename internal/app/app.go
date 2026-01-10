// internal/app/app.go
package app

import (
	"airops/internal/application/usecase"
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
	// Repositories
	flightsRepo := repositories.NewFlightsRepo(pool)
	passengersRepo := repositories.NewPassengersRepo(pool)
	statsRoutesRepo := repositories.NewStatsRoutesRepo(pool)
	healthRepo := repositories.NewHealthRepo(pool)
	bookingsRepo := repositories.NewBookingsRepo(pool)
	seatsRepo := repositories.NewSeatsRepo(pool)
	ticketsRepo := repositories.NewTicketsRepo(pool)
	airportsRepo := repositories.NewAirportsRepo(pool)
	airplanesRepo := repositories.NewAirplanesRepo(pool) // ✨ NEW

	// Services
	flightsService := usecase.NewFlightsService(flightsRepo, passengersRepo)
	passengersService := usecase.NewPassengersService(passengersRepo)
	statsService := usecase.NewStatsRoutesService(statsRoutesRepo)
	healthService := usecase.NewHealthService(healthRepo)
	bookingService := usecase.NewBookingService(bookingsRepo, flightsRepo, seatsRepo, ticketsRepo)
	searchService := usecase.NewSearchService(flightsRepo, seatsRepo)
	airportsService := usecase.NewAirportsService(airportsRepo)
	airplanesService := usecase.NewAirplanesService(airplanesRepo) // ✨ NEW

	// Handlers
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

	// Server
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		server: server,
		pool:   pool,
	}
}

func (a *App) Run() error {
	// Запуск сервера в горутине
	go func() {
		log.Printf("🚀 Starting server on %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	log.Println("✅ Server stopped gracefully")
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	// Закрываем HTTP сервер
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	// Закрываем пул соединений с БД
	a.pool.Close()

	return nil
}
