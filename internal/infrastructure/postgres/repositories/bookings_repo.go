package repositories

import (
	"airops/internal/domain/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingsRepo struct {
	pool *pgxpool.Pool
}

func NewBookingsRepo(pool *pgxpool.Pool) *BookingsRepo {
	return &BookingsRepo{pool: pool}
}

// создает бронирование ВНУТРИ ТРАНЗАКЦИИ
func (r *BookingsRepo) Create(ctx context.Context, tx pgx.Tx, booking *models.Booking) error {
	query := `
		INSERT INTO bookings.bookings 
		(book_ref, book_date, total_amount)
		VALUES ($1, $2, $3)
	`

	_, err := tx.Exec(ctx, query,
		booking.BookRef,
		booking.BookDate,
		booking.TotalAmount,
	)
	if err != nil {
		return fmt.Errorf("create booking: %w", err)
	}
	return nil
}

// получает бронирование по номеру
func (r *BookingsRepo) GetByRef(ctx context.Context, bookRef string) (*models.Booking, error) {
	query := `
		SELECT book_ref, book_date, total_amount
		FROM bookings.bookings
		WHERE book_ref = $1
	`

	var booking models.Booking
	err := r.pool.QueryRow(ctx, query, bookRef).Scan(
		&booking.BookRef,
		&booking.BookDate,
		&booking.TotalAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("get booking by ref: %w", err)
	}

	return &booking, nil
}

// получает все бронирования пассажира
func (r *BookingsRepo) GetByPassenger(ctx context.Context, passengerID string) ([]models.Booking, error) {
	query := `
		SELECT DISTINCT b.book_ref, b.book_date, b.total_amount
		FROM bookings.bookings b
		JOIN bookings.tickets t ON t.book_ref = b.book_ref
		WHERE t.passenger_id = $1
		ORDER BY b.book_date DESC
	`

	rows, err := r.pool.Query(ctx, query, passengerID)
	if err != nil {
		return nil, fmt.Errorf("get bookings by passenger: %w", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.BookRef,
			&booking.BookDate,
			&booking.TotalAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return bookings, nil
}

// отменяет бронирование ВНУТРИ ТРАНЗАКЦИИ
func (r *BookingsRepo) Cancel(ctx context.Context, tx pgx.Tx, bookRef string) error {
	// 1. Удаляем boarding passes
	deleteBoardingPasses := `
		DELETE FROM bookings.boarding_passes bp
		USING bookings.tickets t
		WHERE bp.ticket_no = t.ticket_no
		  AND t.book_ref = $1
	`
	if _, err := tx.Exec(ctx, deleteBoardingPasses, bookRef); err != nil {
		return fmt.Errorf("delete boarding passes: %w", err)
	}

	// 2. Удаляем сегменты
	deleteSegments := `
		DELETE FROM bookings.segments s
		USING bookings.tickets t
		WHERE s.ticket_no = t.ticket_no
		  AND t.book_ref = $1
	`
	if _, err := tx.Exec(ctx, deleteSegments, bookRef); err != nil {
		return fmt.Errorf("delete segments: %w", err)
	}

	// 3. Удаляем билеты
	deleteTickets := `
		DELETE FROM bookings.tickets
		WHERE book_ref = $1
	`
	if _, err := tx.Exec(ctx, deleteTickets, bookRef); err != nil {
		return fmt.Errorf("delete tickets: %w", err)
	}

	// 4. Удаляем бронирование
	deleteBooking := `
		DELETE FROM bookings.bookings
		WHERE book_ref = $1
	`
	result, err := tx.Exec(ctx, deleteBooking, bookRef)
	if err != nil {
		return fmt.Errorf("delete booking: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

// получает бронирование с билетами и рейсами
func (r *BookingsRepo) GetWithDetails(ctx context.Context, bookRef string) (*models.BookingDetails, error) {
	// Получаем само бронирование
	booking, err := r.GetByRef(ctx, bookRef)
	if err != nil {
		return nil, err
	}

	// Получаем билеты
	ticketsQuery := `
		SELECT t.ticket_no, t.passenger_id, t.passenger_name,
		       s.flight_id, s.fare_conditions, s.price,
		       bp.seat_no, bp.boarding_no,
		       f.route_no, f.status, f.scheduled_departure, f.scheduled_arrival
		FROM bookings.tickets t
		JOIN bookings.segments s ON s.ticket_no = t.ticket_no
		LEFT JOIN bookings.boarding_passes bp ON bp.ticket_no = t.ticket_no AND bp.flight_id = s.flight_id
		LEFT JOIN bookings.flights f ON f.flight_id = s.flight_id
		WHERE t.book_ref = $1
		ORDER BY t.ticket_no, s.flight_id
	`

	rows, err := r.pool.Query(ctx, ticketsQuery, bookRef)
	if err != nil {
		return nil, fmt.Errorf("get booking details: %w", err)
	}
	defer rows.Close()

	details := &models.BookingDetails{
		BookRef:     booking.BookRef,
		BookDate:    booking.BookDate,
		TotalAmount: booking.TotalAmount,
		Tickets:     []models.TicketDetails{},
	}

	ticketsMap := make(map[string]*models.TicketDetails)

	for rows.Next() {
		var (
			ticketNo       string
			passengerID    string
			passengerName  string
			flightID       int64
			fareConditions string
			price          string
			seatNo         *string
			boardingNo     *int
			routeNo        string
			status         string
			schedDep       time.Time
			schedArr       time.Time
		)

		err := rows.Scan(
			&ticketNo, &passengerID, &passengerName,
			&flightID, &fareConditions, &price,
			&seatNo, &boardingNo,
			&routeNo, &status, &schedDep, &schedArr,
		)
		if err != nil {
			return nil, fmt.Errorf("scan ticket details: %w", err)
		}

		// Если билета еще нет в map - создаем
		if _, ok := ticketsMap[ticketNo]; !ok {
			ticketsMap[ticketNo] = &models.TicketDetails{
				TicketNo:      ticketNo,
				PassengerID:   passengerID,
				PassengerName: passengerName,
				Flights:       []models.FlightInfo{},
			}
		}

		// Добавляем информацию о рейсе
		flightInfo := models.FlightInfo{
			FlightID:           flightID,
			RouteNo:            routeNo,
			Status:             status,
			FareConditions:     fareConditions,
			ScheduledDeparture: schedDep,
			ScheduledArrival:   schedArr,
		}
		if seatNo != nil {
			flightInfo.SeatNo = *seatNo
		}
		if boardingNo != nil {
			flightInfo.BoardingNo = *boardingNo
		}

		ticketsMap[ticketNo].Flights = append(ticketsMap[ticketNo].Flights, flightInfo)
	}

	// Конвертируем map в slice
	for _, ticket := range ticketsMap {
		details.Tickets = append(details.Tickets, *ticket)
	}

	return details, nil
}

// BeginTx начинает транзакцию (для использования в сервисах)
func (r *BookingsRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}
