package usecase

import (
	"airops/internal/domain/models"
	repositories2 "airops/internal/infrastructure/postgres/repositories"
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type BookingService struct {
	bookingsRepo *repositories2.BookingsRepo
	flightsRepo  *repositories2.FlightsRepo
	seatsRepo    *repositories2.SeatsRepo
	ticketsRepo  *repositories2.TicketsRepo
}

func NewBookingService(
	bookingsRepo *repositories2.BookingsRepo,
	flightsRepo *repositories2.FlightsRepo,
	seatsRepo *repositories2.SeatsRepo,
	ticketsRepo *repositories2.TicketsRepo,
) *BookingService {
	return &BookingService{
		bookingsRepo: bookingsRepo,
		flightsRepo:  flightsRepo,
		seatsRepo:    seatsRepo,
		ticketsRepo:  ticketsRepo,
	}
}

// создает новое бронирование
func (s *BookingService) Create(ctx context.Context, req models.BookingRequest) (*models.BookingDetails, error) {
	// Проверяем доступность рейса
	flight, err := s.flightsRepo.GetByID(ctx, req.FlightID)
	if err != nil {
		return nil, fmt.Errorf("flight not found: %w", err)
	}

	// Проверяем статус рейса
	if flight.Status != "Scheduled" && flight.Status != "On Time" {
		return nil, fmt.Errorf("flight is not available for booking (status: %s)", flight.Status)
	}

	// Проверяем что рейс еще не вылетел
	if flight.ScheduledDeparture.Before(time.Now()) {
		return nil, fmt.Errorf("flight has already departed")
	}

	// Проверяем доступность мест
	if len(req.Seats) == 0 {
		return nil, fmt.Errorf("no seats selected")
	}

	for _, seatNo := range req.Seats {
		available, err := s.seatsRepo.IsSeatAvailable(ctx, req.FlightID, seatNo)
		if err != nil {
			return nil, fmt.Errorf("check seat availability: %w", err)
		}
		if !available {
			return nil, fmt.Errorf("seat %s is not available", seatNo)
		}
	}

	// Рассчитываем стоимость
	totalAmount := s.calculatePrice(req.Seats, req.FareClass)

	// Начинаем транзакцию
	tx, err := s.bookingsRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Генерируем уникальные идентификаторы
	bookRef := generateBookRef()
	ticketNo := generateTicketNo()

	// Создаем бронирование
	booking := &models.Booking{
		BookRef:     bookRef,
		BookDate:    time.Now(),
		TotalAmount: totalAmount,
	}

	err = s.bookingsRepo.Create(ctx, tx, booking)
	if err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	// Создаем билет
	ticket := &models.Ticket{
		TicketNo:      ticketNo,
		BookRef:       bookRef,
		PassengerID:   req.PassengerID,
		PassengerName: req.PassengerName,
		Outbound:      true,
	}

	err = s.ticketsRepo.Create(ctx, tx, ticket)
	if err != nil {
		return nil, fmt.Errorf("create ticket: %w", err)
	}

	// Создаем сегмент (привязка билета к рейсу)
	segment := &models.TicketSegment{
		TicketNo:       ticketNo,
		FlightID:       req.FlightID,
		FareConditions: req.FareClass,
		Price:          totalAmount,
	}

	err = s.ticketsRepo.CreateSegment(ctx, tx, segment)
	if err != nil {
		return nil, fmt.Errorf("create segment: %w", err)
	}

	// Резервируем места (создаем boarding passes)
	for _, seatNo := range req.Seats {
		err = s.seatsRepo.Reserve(ctx, tx, req.FlightID, ticketNo, seatNo)
		if err != nil {
			return nil, fmt.Errorf("reserve seat: %w", err)
		}
	}

	// Коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	// Получаем полные детали бронирования
	details, err := s.bookingsRepo.GetWithDetails(ctx, bookRef)
	if err != nil {
		return nil, fmt.Errorf("get booking details: %w", err)
	}

	return details, nil
}

// получает бронирование по номеру
func (s *BookingService) GetByRef(ctx context.Context, bookRef string) (*models.BookingDetails, error) {
	details, err := s.bookingsRepo.GetWithDetails(ctx, bookRef)
	if err != nil {
		return nil, fmt.Errorf("get booking: %w", err)
	}

	return details, nil
}

// получает все бронирования пассажира
func (s *BookingService) GetByPassenger(ctx context.Context, passengerID string) ([]models.Booking, error) {
	bookings, err := s.bookingsRepo.GetByPassenger(ctx, passengerID)
	if err != nil {
		return nil, fmt.Errorf("get bookings by passenger: %w", err)
	}

	return bookings, nil
}

// отменяет бронирование
func (s *BookingService) Cancel(ctx context.Context, bookRef string) error {
	// Начинаем транзакцию
	tx, err := s.bookingsRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Проверяем что бронирование существует (можно и без этого, но ок)
	_, err = s.bookingsRepo.GetByRef(ctx, bookRef)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	// Удаляем всё каскадно в рамках tx
	if err := s.bookingsRepo.Cancel(ctx, tx, bookRef); err != nil {
		return fmt.Errorf("cancel booking: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// рассчитывает стоимость билетов
func (s *BookingService) calculatePrice(seats []string, fareClass string) decimal.Decimal {
	// Базовые цены по классам обслуживания
	basePrices := map[string]float64{
		"Economy":  5000.0,
		"Comfort":  10000.0,
		"Business": 20000.0,
	}

	basePrice, ok := basePrices[fareClass]
	if !ok {
		basePrice = basePrices["Economy"] // по умолчанию
	}

	// Цена = базовая цена * количество мест
	totalPrice := basePrice * float64(len(seats))

	return decimal.NewFromFloat(totalPrice)
}

// генерирует уникальный 6-символьный код бронирования
func generateBookRef() string {
	b := make([]byte, 4)
	rand.Read(b)
	// Кодируем в base32 и берем первые 6 символов
	encoded := base32.StdEncoding.EncodeToString(b)
	return strings.ToUpper(encoded[:6])
}

// генерирует уникальный номер билета (13 цифр)
func generateTicketNo() string {
	// Формат: YYYYMMDDHHMMSS (timestamp)
	return time.Now().Format("20060102150405") + fmt.Sprintf("%03d", time.Now().Nanosecond()/1000000)
}
