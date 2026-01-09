// stat_service.go
package usecase

import (
	"airops/internal/domain"
	"airops/internal/domain/models"
	"context"
	"time"
)

type StatsRoutesService struct {
	stats domain.StatsRoutesRepo
}

func NewStatsRoutesService(s domain.StatsRoutesRepo) *StatsRoutesService {
	return &StatsRoutesService{stats: s}
}

func (s *StatsRoutesService) TopRoutes(
	ctx context.Context,
	from, to time.Time,
	limit int,
) ([]models.RouteStat, error) {
	if limit <= 0 {
		limit = 10
	}

	items, err := s.stats.TopRoutes(ctx, from, to, limit)
	if err != nil {
		return nil, MapStoreErr(err)
	}

	return items, nil
}
