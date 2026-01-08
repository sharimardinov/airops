package app

import (
	nethttp "net/http"

	"airops/internal/infra/db/pg/repo"
	transporthttp "airops/internal/transport/http"
	"airops/internal/transport/http/handlers"
	"airops/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Handler nethttp.Handler
}

func New(pool *pgxpool.Pool) *App {
	// repos (имена 1 в 1 как на твоих скринах)
	flightsRepo := repo.NewFlightsRepo(pool)
	passengersRepo := repo.NewPassengersRepo(pool)
	statsRepo := repo.NewStatsRoutesRepo(pool)

	// usecases
	// ВАЖНО: FlightsService теперь будет принимать passengersRepo
	flightsUC := usecase.NewFlightsService(flightsRepo, passengersRepo)
	passengersUC := usecase.NewPassengersService(passengersRepo)
	statsUC := usecase.NewStatsRoutesService(statsRepo)

	// handlers
	h := handlers.New(flightsUC, passengersUC, statsUC)

	// router (package http => импорт алиасом transporthttp)
	r := transporthttp.New(h)

	return &App{Handler: r}
}
