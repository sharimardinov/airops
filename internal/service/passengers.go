package service

import (
	"airops/internal/store"
	"context"
)

type PassengersService struct {
	store *store.PassengersStore
}

func NewPassengersService(s *store.PassengersStore) *PassengersService {
	return &PassengersService{store: s}
}

func (s *PassengersService) ListByFlightID(
	ctx context.Context,
	flightID int64,
	limit, offset int,
) ([]store.PassengerRow, error) {
	rows, err := s.store.ListByFlightID(ctx, flightID, limit, offset)
	if err != nil {
		return nil, mapStoreErr(err)
	}
	return rows, nil
}
