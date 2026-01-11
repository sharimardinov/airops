package repositories

import (
	"airops/internal/domain/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketsRepo struct {
	pool *pgxpool.Pool
}

func NewTicketsRepo(pool *pgxpool.Pool) *TicketsRepo {
	return &TicketsRepo{pool: pool}
}

// создает билет в транзакции
func (r *TicketsRepo) Create(ctx context.Context, tx pgx.Tx, ticket *models.Ticket) error {
	query := `
		INSERT INTO bookings.tickets 
		(ticket_no, book_ref, passenger_id, passenger_name, outbound)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := tx.Exec(ctx, query,
		ticket.TicketNo,
		ticket.BookRef,
		ticket.PassengerID,
		ticket.PassengerName,
		ticket.Outbound,
	)
	if err != nil {
		return fmt.Errorf("create ticket: %w", err)
	}

	return nil
}

// создает сегмент билета (привязка к рейсу)
func (r *TicketsRepo) CreateSegment(ctx context.Context, tx pgx.Tx, segment *models.TicketSegment) error {
	query := `
		INSERT INTO bookings.segments 
		(ticket_no, flight_id, fare_conditions, price)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.Exec(ctx, query,
		segment.TicketNo,
		segment.FlightID,
		segment.FareConditions,
		segment.Price,
	)
	if err != nil {
		return fmt.Errorf("create segment: %w", err)
	}

	return nil
}

// получает билет по номеру
func (r *TicketsRepo) GetByNumber(ctx context.Context, ticketNo string) (*models.Ticket, error) {
	query := `
		SELECT ticket_no, book_ref, passenger_id, passenger_name, outbound
		FROM bookings.tickets
		WHERE ticket_no = $1
	`

	var ticket models.Ticket
	err := r.pool.QueryRow(ctx, query, ticketNo).Scan(
		&ticket.TicketNo,
		&ticket.BookRef,
		&ticket.PassengerID,
		&ticket.PassengerName,
		&ticket.Outbound,
	)
	if err != nil {
		return nil, fmt.Errorf("get ticket by number: %w", err)
	}

	return &ticket, nil
}

// получает все билеты по номеру бронирования
func (r *TicketsRepo) GetByBooking(ctx context.Context, bookRef string) ([]models.Ticket, error) {
	query := `
		SELECT ticket_no, book_ref, passenger_id, passenger_name, outbound
		FROM bookings.tickets
		WHERE book_ref = $1
		ORDER BY ticket_no
	`

	rows, err := r.pool.Query(ctx, query, bookRef)
	if err != nil {
		return nil, fmt.Errorf("get tickets by booking: %w", err)
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var ticket models.Ticket
		err := rows.Scan(
			&ticket.TicketNo,
			&ticket.BookRef,
			&ticket.PassengerID,
			&ticket.PassengerName,
			&ticket.Outbound,
		)
		if err != nil {
			return nil, fmt.Errorf("scan ticket: %w", err)
		}
		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tickets, nil
}

// получает все сегменты (рейсы) для билета
func (r *TicketsRepo) GetSegmentsByTicket(ctx context.Context, ticketNo string) ([]models.TicketSegment, error) {
	query := `
		SELECT ticket_no, flight_id, fare_conditions, price
		FROM bookings.segments
		WHERE ticket_no = $1
		ORDER BY flight_id
	`

	rows, err := r.pool.Query(ctx, query, ticketNo)
	if err != nil {
		return nil, fmt.Errorf("get segments: %w", err)
	}
	defer rows.Close()

	var segments []models.TicketSegment
	for rows.Next() {
		var segment models.TicketSegment
		err := rows.Scan(
			&segment.TicketNo,
			&segment.FlightID,
			&segment.FareConditions,
			&segment.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("scan segment: %w", err)
		}
		segments = append(segments, segment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return segments, nil
}

// получает все билеты пассажира
func (r *TicketsRepo) GetByPassenger(ctx context.Context, passengerID string) ([]models.Ticket, error) {
	query := `
		SELECT ticket_no, book_ref, passenger_id, passenger_name, outbound
		FROM bookings.tickets
		WHERE passenger_id = $1
		ORDER BY ticket_no DESC
	`

	rows, err := r.pool.Query(ctx, query, passengerID)
	if err != nil {
		return nil, fmt.Errorf("get tickets by passenger: %w", err)
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var ticket models.Ticket
		err := rows.Scan(
			&ticket.TicketNo,
			&ticket.BookRef,
			&ticket.PassengerID,
			&ticket.PassengerName,
			&ticket.Outbound,
		)
		if err != nil {
			return nil, fmt.Errorf("scan ticket: %w", err)
		}
		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tickets, nil
}

// удаляет билет в транзакции
func (r *TicketsRepo) Delete(ctx context.Context, tx pgx.Tx, ticketNo string) error {
	// Сначала удаляем сегменты
	deleteSegments := `
		DELETE FROM bookings.segments
		WHERE ticket_no = $1
	`
	_, err := tx.Exec(ctx, deleteSegments, ticketNo)
	if err != nil {
		return fmt.Errorf("delete segments: %w", err)
	}

	// Потом сам билет
	deleteTicket := `
		DELETE FROM bookings.tickets
		WHERE ticket_no = $1
	`
	result, err := tx.Exec(ctx, deleteTicket, ticketNo)
	if err != nil {
		return fmt.Errorf("delete ticket: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("ticket not found")
	}

	return nil
}
