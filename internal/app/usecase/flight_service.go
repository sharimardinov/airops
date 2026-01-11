package usecase

import (
	"airops/internal/domain"
	"airops/internal/domain/models"
	"airops/internal/infrastructure/postgres/repositories"
	"context"
	"time"
)

type FlightsService struct {
	flightsRepo    domain.FlightsRepo
	passengersRepo domain.PassengersRepo
}

func NewFlightsService(flightsRepo *repositories.FlightsRepo, passengersRepo domain.PassengersRepo) *FlightsService {
	return &FlightsService{
		flightsRepo:    flightsRepo,
		passengersRepo: passengersRepo,
	}
}

// получает детальную информацию о рейсе
func (s *FlightsService) GetByID(ctx context.Context, id int64) (models.FlightDetails, error) {
	flight, err := s.flightsRepo.GetByID(ctx, id)
	if err != nil {
		return models.FlightDetails{}, MapStoreErr(err)
	}

	passengers, err := s.passengersRepo.ListByFlightID(ctx, id, 0, 1000)
	if err != nil {
		return models.FlightDetails{}, MapStoreErr(err)
	}

	details := models.FlightDetails{
		FlightID:           flight.FlightID,
		RouteNo:            flight.RouteNo,
		Status:             flight.Status,
		ScheduledDeparture: flight.ScheduledDeparture,
		ScheduledArrival:   flight.ScheduledArrival,
		ActualDeparture:    flight.ActualDeparture,
		ActualArrival:      flight.ActualArrival,
		Passengers:         passengers,
	}

	return details, nil
}

// возвращает список рейсов
func (s *FlightsService) List(ctx context.Context, date time.Time, limit, offset int) ([]models.Flight, error) {
	// !!!date приходит как time.Parse -> UTC
	loc := time.Local

	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	to := from.Add(24 * time.Hour)

	return s.flightsRepo.List(ctx, from, to, limit, offset)
}

func (s *FlightsService) listByDate(ctx context.Context, date time.Time, limit int) ([]models.Flight, error) {
	dateEnd := date.Add(24 * time.Hour)
	return s.flightsRepo.List(ctx, date, dateEnd, 0, limit)
}
