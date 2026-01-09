package usecase

import "context"

type HealthPinger interface {
	Ping(ctx context.Context) error
}

type HealthService struct {
	db HealthPinger
}

func NewHealthService(db HealthPinger) *HealthService {
	return &HealthService{db: db}
}

func (s *HealthService) Ready(ctx context.Context) error {
	return s.db.Ping(ctx)
}
