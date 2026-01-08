package usecase

import (
	"airops/internal/domain"
	"airops/internal/domain/models"
	"context"
	"time"
)

type FlightsService struct {
	flights    domain.FlightsRepo
	passengers domain.PassengersRepo
}

func NewFlightsService(f domain.FlightsRepo, p domain.PassengersRepo) *FlightsService {
	return &FlightsService{
		flights:    f,
		passengers: p,
	}
}

func (s *FlightsService) List(
	ctx context.Context,
	from, to time.Time,
	limit, offset int,
) ([]models.Flight, error) {
	items, err := s.flights.List(ctx, from, to, limit, offset)
	if err != nil {
		return nil, MapStoreErr(err)
	}
	return items, nil
}

// Возвращаем FlightDetails без отдельного repo.
// Пассажиров берём "всё" (лимит большой). Пагинация у тебя отдельно в /flights/{id}/passengers.
func (s *FlightsService) GetByID(ctx context.Context, id int64) (models.FlightDetails, error) {
	if id <= 0 {
		return models.FlightDetails{}, domain.ErrInvalidArgument
	}

	flight, err := s.flights.GetByID(ctx, id)
	if err != nil {
		return models.FlightDetails{}, MapStoreErr(err)
	}

	// если боишься больших рейсов — поставь лимит 1000/5000, но сейчас так проще
	passengers, err := s.passengers.ListByFlightID(ctx, id, 10000, 0)
	if err != nil {
		return models.FlightDetails{}, MapStoreErr(err)
	}

	return models.FlightDetails{
		Flight:     flight,
		Passengers: passengers,
	}, nil
}
