// ================================================
// internal/application/usecase/flight_service.go
// ================================================
package usecase

import (
	"airops/internal/domain"
	"airops/internal/domain/models"
	"airops/internal/infrastructure/postgres/repositories"
	"context"
	"time"
)

type FlightsService struct {
	flightsRepo    domain.FlightsRepo
	passengersRepo domain.PassengersRepo
}

func NewFlightsService(flightsRepo *repositories.FlightsRepo, passengersRepo domain.PassengersRepo) *FlightsService {
	return &FlightsService{
		flightsRepo:    flightsRepo,
		passengersRepo: passengersRepo,
	}
}

// GetByID получает детальную информацию о рейсе
func (s *FlightsService) GetByID(ctx context.Context, id int64) (models.FlightDetails, error) {
	// Получаем базовую информацию о рейсе
	flight, err := s.flightsRepo.GetByID(ctx, id)
	if err != nil {
		return models.FlightDetails{}, MapStoreErr(err)
	}

	// ✅ ИСПРАВЛЕНО: добавлены недостающие параметры
	// Сигнатура: ListByFlightID(ctx, flightID, offset, limit)
	passengers, err := s.passengersRepo.ListByFlightID(ctx, id, 0, 1000)
	if err != nil {
		return models.FlightDetails{}, MapStoreErr(err)
	}

	// Создаем FlightDetails
	details := models.FlightDetails{
		FlightID:           flight.FlightID,
		RouteNo:            flight.RouteNo,
		Status:             flight.Status,
		ScheduledDeparture: flight.ScheduledDeparture,
		ScheduledArrival:   flight.ScheduledArrival,
		ActualDeparture:    flight.ActualDeparture,
		ActualArrival:      flight.ActualArrival,
		Passengers:         passengers,
	}

	return details, nil
}

// List возвращает список рейсов
func (s *FlightsService) List(ctx context.Context, date time.Time, limit int) ([]models.Flight, error) {
	// Сигнатура: List(ctx, dateFrom, dateTo, offset, limit)

	// Вариант 1: Если нужны рейсы только за один день
	dateEnd := date.Add(24 * time.Hour)
	flights, err := s.flightsRepo.List(ctx, date, dateEnd, 0, limit)

	if err != nil {
		return nil, MapStoreErr(err)
	}
	return flights, nil
}

// Вспомогательный метод (если нужен)
func (s *FlightsService) listByDate(ctx context.Context, date time.Time, limit int) ([]models.Flight, error) {
	dateEnd := date.Add(24 * time.Hour)
	return s.flightsRepo.List(ctx, date, dateEnd, 0, limit)
}
