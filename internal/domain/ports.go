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

type BookingRepo interface {
	Create(ctx context.Context, booking *models.Booking) error
	GetByRef(ctx context.Context, bookRef string) (*models.Booking, error)
	GetByPassenger(ctx context.Context, passengerID string) ([]models.Booking, error)
	Cancel(ctx context.Context, bookRef string) error
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

type AirportRepo interface {
	GetByCode(ctx context.Context, code string) (*models.Airport, error)
	List(ctx context.Context) ([]models.Airport, error)
	SearchByCity(ctx context.Context, city string) ([]models.Airport, error)
}

type SeatRepo interface {
	GetByAirplane(ctx context.Context, airplaneCode string) ([]models.Seat, error)
	GetAvailableByFlight(ctx context.Context, flightID int, fareClass string) ([]models.Seat, error)
	Reserve(ctx context.Context, flightID int, seatNo string) error
}
