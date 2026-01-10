// internal/app/usecase/search_service.go
package usecase

import (
	"airops/internal/domain/models"
	"airops/internal/infrastructure/postgres/repositories"
	"context"
	"fmt"
)

type SearchService struct {
	flightsRepo *repositories.FlightsRepo
	seatsRepo   *repositories.SeatsRepo
}

func NewSearchService(
	flightsRepo *repositories.FlightsRepo,
	seatsRepo *repositories.SeatsRepo,
) *SearchService {
	return &SearchService{
		flightsRepo: flightsRepo,
		seatsRepo:   seatsRepo,
	}
}

// SearchFlights ищет рейсы по параметрам
func (s *SearchService) SearchFlights(ctx context.Context, params models.FlightSearchParams) ([]models.FlightSearchResult, error) {
	// Валидация
	if params.DepartureAirport == "" || params.ArrivalAirport == "" {
		return nil, fmt.Errorf("departure and arrival airports are required")
	}
	if params.DepartureDate.IsZero() {
		return nil, fmt.Errorf("departure date is required")
	}
	if params.Passengers <= 0 {
		params.Passengers = 1
	}

	// ✅ ИСПОЛЬЗУЕМ МЕТОД Search (если он есть)
	// Если метода Search еще нет, используй временное решение ниже
	flights, err := s.flightsRepo.Search(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("search flights: %w", err)
	}

	// Обогащаем результаты
	results := make([]models.FlightSearchResult, 0, len(flights))
	for _, flight := range flights {
		// Проверяем доступность мест
		availableCount, err := s.seatsRepo.GetAvailableCount(ctx, flight.FlightID, params.FareClass)
		if err != nil {
			continue
		}

		// Фильтруем по количеству пассажиров
		if availableCount < params.Passengers {
			continue
		}

		// Рассчитываем цену
		price := s.calculatePrice(params.Passengers, params.FareClass)

		// ✅ ПРАВИЛЬНО: встраиваем FlightDetails
		result := models.FlightSearchResult{
			FlightDetails:  flight,
			AvailableSeats: availableCount,
			Price:          price,
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *SearchService) calculatePrice(passengers int, fareClass string) float64 {
	basePrices := map[string]float64{
		"Economy":  5000.0,
		"Comfort":  10000.0,
		"Business": 20000.0,
	}

	basePrice, ok := basePrices[fareClass]
	if !ok {
		basePrice = basePrices["Economy"]
	}

	return basePrice * float64(passengers)
}
