package service

import (
	"airops/internal/store"
	"context"
	"time"
)

type FlightsService struct {
	flights *store.FlightsStore
	details *store.FlightDetailsStore
}

func NewFlightsService(f *store.FlightsStore, d *store.FlightDetailsStore) *FlightsService {
	return &FlightsService{
		flights: f,
		details: d,
	}
}

func (s *FlightsService) List(
	ctx context.Context,
	from, to time.Time,
	limit, offset int,
) ([]store.Flight, error) {
	flights, err := s.flights.List(ctx, from, to, limit, offset)
	if err != nil {
		return nil, mapStoreErr(err)
	}
	return flights, nil
}

func (s *FlightsService) GetByID(ctx context.Context, id int64) (store.FlightDetails, error) {
	details, err := s.details.GetByID(ctx, id)
	if err != nil {
		return store.FlightDetails{}, mapStoreErr(err)
	}
	return details, nil
}
