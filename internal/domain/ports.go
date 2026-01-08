package domain

import (
	"airops/internal/domain/models"
	"context"
)
import "time"

type FlightsRepo interface {
	List(ctx context.Context, from, to time.Time, limit, offset int) ([]models.Flight, error)
	GetByID(ctx context.Context, id int64) (models.Flight, error)
}

type FlightDetailsRepo interface {
	Get(ctx context.Context, id int64) (models.FlightDetails, error)
}

type PassengersRepo interface {
	ListByFlightID(ctx context.Context, flightID int64, limit, offset int) ([]models.FlightPassenger, error)
}

type StatsRoutesRepo interface {
	TopRoutes(ctx context.Context, from, to time.Time, limit int) ([]models.RouteStat, error)
}
