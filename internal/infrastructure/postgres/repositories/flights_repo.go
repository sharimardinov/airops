package repositories

import (
	"airops/internal/domain/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FlightsRepo struct {
	pool *pgxpool.Pool
}

func NewFlightsRepo(pool *pgxpool.Pool) *FlightsRepo {
	return &FlightsRepo{pool: pool}
}

// GetByID получает рейс по ID (только базовые поля из таблицы flights)
func (r *FlightsRepo) GetByID(ctx context.Context, id int64) (models.Flight, error) {
	query := `
		SELECT 
			flight_id,
			route_no,
			status,
			scheduled_departure,
			scheduled_arrival,
			actual_departure,
			actual_arrival
		FROM bookings.flights
		WHERE flight_id = $1
	`

	var f models.Flight
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&f.FlightID, // ✅ было &f.ID
		&f.RouteNo,
		&f.Status,
		&f.ScheduledDeparture,
		&f.ScheduledArrival,
		&f.ActualDeparture,
		&f.ActualArrival,
		// ✅ НЕ сканируем DepartureAirport и ArrivalAirport - их нет в базовом Flight
	)
	if err != nil {
		return models.Flight{}, fmt.Errorf("get flight by id: %w", err)
	}

	return f, nil
}

func (r *FlightsRepo) List(
	ctx context.Context,
	from time.Time,
	to time.Time,
	limit int,
	offset int,
) ([]models.Flight, error) {

	const q = `
		SELECT 
			flight_id,
			route_no,
			status,
			scheduled_departure,
			scheduled_arrival,
			actual_departure,
			actual_arrival
		FROM bookings.flights
		WHERE scheduled_departure >= $1
		  AND scheduled_departure < $2
		ORDER BY scheduled_departure  -- ✅ ИСПРАВЛЕНО: было flight_no
		OFFSET $3 LIMIT $4
`

	rows, err := r.pool.Query(ctx, q, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list flights: %w", err)
	}
	defer rows.Close()

	out := make([]models.Flight, 0, limit)
	for rows.Next() {
		var f models.Flight
		if err := rows.Scan(
			&f.FlightID,
			&f.RouteNo,
			&f.Status,
			&f.ScheduledDeparture,
			&f.ScheduledArrival,
			&f.ActualDeparture,
			&f.ActualArrival,
		); err != nil {
			return nil, fmt.Errorf("scan flight: %w", err)
		}
		out = append(out, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return out, nil
}

// ✅ ДОБАВЬ НОВЫЙ МЕТОД Search для полноценного поиска с JOIN
func (r *FlightsRepo) Search(ctx context.Context, params models.FlightSearchParams) ([]models.FlightDetails, error) {
	query := `
		SELECT 
			f.flight_id,
			f.route_no,
			f.status,
			f.scheduled_departure,
			f.scheduled_arrival,
			f.actual_departure,
			f.actual_arrival,
			r.departure_airport,
			r.arrival_airport,
			r.airplane_code,
			dep.airport_name as departure_airport_name,
			dep.city as departure_city,
			arr.airport_name as arrival_airport_name,
			arr.city as arrival_city,
			COALESCE(a.model, '') as airplane_model,
			COALESCE(r.duration, INTERVAL '0') as duration
		FROM bookings.flights f
		JOIN bookings.routes r ON r.route_no = f.route_no
		JOIN bookings.airports dep ON dep.airport_code = r.departure_airport
		JOIN bookings.airports arr ON arr.airport_code = r.arrival_airport
		LEFT JOIN bookings.airplanes a ON a.airplane_code = r.airplane_code
		WHERE r.departure_airport = $1
		  AND r.arrival_airport = $2
		  AND f.scheduled_departure::date = $3
		  AND f.status IN ('Scheduled', 'On Time')
		ORDER BY f.scheduled_departure
	`

	rows, err := r.pool.Query(ctx, query,
		params.DepartureAirport,
		params.ArrivalAirport,
		params.DepartureDate,
	)
	if err != nil {
		return nil, fmt.Errorf("search flights: %w", err)
	}
	defer rows.Close()

	var flights []models.FlightDetails
	for rows.Next() {
		var fd models.FlightDetails

		err := rows.Scan(
			&fd.FlightID,
			&fd.RouteNo,
			&fd.Status,
			&fd.ScheduledDeparture,
			&fd.ScheduledArrival,
			&fd.ActualDeparture,
			&fd.ActualArrival,
			&fd.DepartureAirport,
			&fd.ArrivalAirport,
			&fd.AirplaneCode,
			&fd.DepartureAirportName,
			&fd.DepartureCity,
			&fd.ArrivalAirportName,
			&fd.ArrivalCity,
			&fd.AirplaneModel,
			&fd.Duration,
		)
		if err != nil {
			return nil, fmt.Errorf("scan flight: %w", err)
		}

		flights = append(flights, fd)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return flights, nil
}
