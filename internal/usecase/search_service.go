// internal/usecase/search_service.go
package usecase

import (
	"airops/internal/domain/models"
	"airops/internal/infra/db/pg/repo"
	"context"
	"fmt"
	"time"
)

type SearchService struct {
	flightsRepo *repo.FlightsRepo
	seatsRepo   *repo.SeatsRepo
}

func NewSearchService(
	flightsRepo *repo.FlightsRepo,
	seatsRepo *repo.SeatsRepo,
) *SearchService {
	return &SearchService{
		flightsRepo: flightsRepo,
		seatsRepo:   seatsRepo,
	}
}

// SearchFlights ищет рейсы по параметрам
func (s *SearchService) SearchFlights(ctx context.Context, params models.FlightSearchParams) ([]models.FlightSearchResult, error) {
	// Валидация параметров
	if params.DepartureAirport == "" || params.ArrivalAirport == "" {
		return nil, fmt.Errorf("departure and arrival airports are required")
	}

	if params.DepartureDate.IsZero() {
		return nil, fmt.Errorf("departure date is required")
	}

	if params.Passengers <= 0 {
		params.Passengers = 1
	}

	// Ищем рейсы
	flights, err := s.searchFlightsByRoute(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("search flights: %w", err)
	}

	// Обогащаем результаты информацией о доступных местах
	results := make([]models.FlightSearchResult, 0, len(flights))
	for _, flight := range flights {
		// Проверяем доступность мест
		availableCount, err := s.seatsRepo.GetAvailableCount(ctx, flight.ID, params.FareClass)
		if err != nil {
			continue // пропускаем рейс при ошибке
		}

		// Фильтруем рейсы по количеству доступных мест
		if availableCount < params.Passengers {
			continue
		}

		// Рассчитываем примерную стоимость
		price := s.calculatePrice(params.Passengers, params.FareClass)

		result := models.FlightSearchResult{
			Flight:         flight,
			AvailableSeats: availableCount,
			Price:          price,
		}

		results = append(results, result)
	}

	return results, nil
}

// searchFlightsByRoute выполняет поиск рейсов в базе
func (s *SearchService) searchFlightsByRoute(ctx context.Context, params models.FlightSearchParams) ([]models.Flight, error) {
	from := params.DepartureDate
	to := params.DepartureDate.Add(24 * time.Hour)

	flights, err := s.flightsRepo.List(ctx, from, to, 100, 0)
	if err != nil {
		return nil, err
	}

	// Фильтр по направлению (у тебя эти поля уже есть в models.Flight)
	filtered := make([]models.Flight, 0, len(flights))
	for _, f := range flights {
		if f.DepartureAirport == params.DepartureAirport && f.ArrivalAirport == params.ArrivalAirport {
			filtered = append(filtered, f)
		}
	}

	return filtered, nil
}

// calculatePrice рассчитывает стоимость для заданного количества пассажиров
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
