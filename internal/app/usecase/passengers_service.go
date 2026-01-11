// passengers_service.go
package usecase

import (
	"airops/internal/domain"
	"airops/internal/domain/models"
	"context"
)

type PassengersService struct {
	passengers domain.PassengersRepo
}

func NewPassengersService(p domain.PassengersRepo) *PassengersService {
	return &PassengersService{passengers: p}
}

func (s *PassengersService) ListByFlightID(
	ctx context.Context,
	flightID int64,
	limit, offset int,
) ([]models.FlightPassenger, error) {
	if flightID <= 0 {
		return nil, domain.ErrInvalidArgument
	}

	items, err := s.passengers.ListByFlightID(ctx, flightID, limit, offset)
	if err != nil {
		return nil, MapStoreErr(err)
	}
	return items, nil
}
