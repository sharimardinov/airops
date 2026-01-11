package usecase

import (
	"airops/internal/domain/models"
	"airops/internal/infrastructure/postgres/repositories"
	"context"
	"fmt"
)

type AirportsService struct {
	airportsRepo *repositories.AirportsRepo
}

func NewAirportsService(airportsRepo *repositories.AirportsRepo) *AirportsService {
	return &AirportsService{
		airportsRepo: airportsRepo,
	}
}

// GetByCode получает аэропорт по коду
func (s *AirportsService) GetByCode(ctx context.Context, code string) (*models.Airport, error) {
	airport, err := s.airportsRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get airport: %w", err)
	}

	return airport, nil
}

// List возвращает список всех аэропортов
func (s *AirportsService) List(ctx context.Context) ([]models.Airport, error) {
	airports, err := s.airportsRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list airports: %w", err)
	}

	return airports, nil
}

// SearchByCity ищет аэропорты по названию города
func (s *AirportsService) SearchByCity(ctx context.Context, city string) ([]models.Airport, error) {
	airports, err := s.airportsRepo.SearchByCity(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("search airports: %w", err)
	}

	return airports, nil
}

// SearchByCountry ищет аэропорты по стране
func (s *AirportsService) SearchByCountry(ctx context.Context, country string) ([]models.Airport, error) {
	airports, err := s.airportsRepo.SearchByCountry(ctx, country)
	if err != nil {
		return nil, fmt.Errorf("search airports: %w", err)
	}

	return airports, nil
}
