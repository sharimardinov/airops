package usecase

import (
	"airops/internal/domain/models"
	"airops/internal/infrastructure/postgres/repositories"
	"context"
	"fmt"
)

type AirplanesService struct {
	airplanesRepo *repositories.AirplanesRepo
}

func NewAirplanesService(airplanesRepo *repositories.AirplanesRepo) *AirplanesService {
	return &AirplanesService{
		airplanesRepo: airplanesRepo,
	}
}

// получает самолет по коду
func (s *AirplanesService) GetByCode(ctx context.Context, code string) (*models.Airplane, error) {
	airplane, err := s.airplanesRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get airplane: %w", err)
	}
	return airplane, nil
}

// возвращает список всех самолетов
func (s *AirplanesService) List(ctx context.Context) ([]models.Airplane, error) {
	airplanes, err := s.airplanesRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list airplanes: %w", err)
	}
	return airplanes, nil
}

// получает самолет с раскладкой мест
func (s *AirplanesService) GetWithSeats(ctx context.Context, code string) (*models.AirplaneWithSeats, error) {
	airplane, err := s.airplanesRepo.GetWithSeats(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get airplane with seats: %w", err)
	}
	return airplane, nil
}

// получает статистику по самолету
func (s *AirplanesService) GetStats(ctx context.Context, code string) (*models.AirplaneStats, error) {
	stats, err := s.airplanesRepo.GetStats(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get airplane stats: %w", err)
	}
	return stats, nil
}
