package service

import (
	"airops/internal/store"
	"context"
)

type StatsRoutesService struct {
	store *store.StatsStore
}

func NewStatsRoutesService(s *store.StatsStore) *StatsRoutesService {
	return &StatsRoutesService{store: s}
}

func (s *StatsRoutesService) Routes(
	ctx context.Context,
	filter store.RoutesStatsFilter,
) ([]store.RouteStat, error) {
	stats, err := s.store.Routes(ctx, filter)
	if err != nil {
		return nil, mapStoreErr(err)
	}
	return stats, nil
}
