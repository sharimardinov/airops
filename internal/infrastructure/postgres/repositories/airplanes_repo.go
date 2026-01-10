// internal/domain/postgres/repositories/airplanes_repo.go
package repositories

import (
	"airops/internal/domain/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AirplanesRepo struct {
	pool *pgxpool.Pool
}

func NewAirplanesRepo(pool *pgxpool.Pool) *AirplanesRepo {
	return &AirplanesRepo{pool: pool}
}

// GetByCode получает самолет по коду
func (r *AirplanesRepo) GetByCode(ctx context.Context, code string) (*models.Airplane, error) {
	query := `
		SELECT airplane_code, model, range, speed
		FROM bookings.airplanes
		WHERE airplane_code = $1
	`

	var airplane models.Airplane
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&airplane.Code,
		&airplane.Model,
		&airplane.Range,
		&airplane.Speed,
	)
	if err != nil {
		return nil, fmt.Errorf("get airplane: %w", err)
	}

	return &airplane, nil
}

// List возвращает все самолеты
func (r *AirplanesRepo) List(ctx context.Context) ([]models.Airplane, error) {
	query := `
		SELECT airplane_code, model, range, speed
		FROM bookings.airplanes
		ORDER BY model
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list airplanes: %w", err)
	}
	defer rows.Close()

	var airplanes []models.Airplane
	for rows.Next() {
		var airplane models.Airplane
		err := rows.Scan(
			&airplane.Code,
			&airplane.Model,
			&airplane.Range,
			&airplane.Speed,
		)
		if err != nil {
			return nil, fmt.Errorf("scan airplane: %w", err)
		}
		airplanes = append(airplanes, airplane)
	}

	return airplanes, rows.Err()
}

// GetWithSeats получает самолет с раскладкой мест
func (r *AirplanesRepo) GetWithSeats(ctx context.Context, code string) (*models.AirplaneWithSeats, error) {
	// Получаем самолет
	airplane, err := r.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Получаем места
	seatsQuery := `
		SELECT seat_no, fare_conditions
		FROM bookings.seats
		WHERE airplane_code = $1
		ORDER BY seat_no
	`

	rows, err := r.pool.Query(ctx, seatsQuery, code)
	if err != nil {
		return nil, fmt.Errorf("get seats: %w", err)
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		seat.AirplaneCode = code
		err := rows.Scan(&seat.SeatNo, &seat.FareConditions)
		if err != nil {
			return nil, fmt.Errorf("scan seat: %w", err)
		}
		seats = append(seats, seat)
	}

	return &models.AirplaneWithSeats{
		Airplane: *airplane,
		Seats:    seats,
	}, nil
}

// GetStats возвращает статистику по самолету
func (r *AirplanesRepo) GetStats(ctx context.Context, code string) (*models.AirplaneStats, error) {
	query := `
WITH seats AS (
  SELECT COUNT(*)::float AS seats_per_plane
  FROM bookings.seats
  WHERE airplane_code = $1
),
flights AS (
  SELECT f.flight_id
  FROM bookings.flights f
  JOIN bookings.routes rt ON rt.route_no = f.route_no
  WHERE rt.airplane_code = $1
    AND f.status = 'Arrived'
),
pax_per_flight AS (
  SELECT bp.flight_id, COUNT(*)::float AS pax
  FROM bookings.boarding_passes bp
  JOIN flights f ON f.flight_id = bp.flight_id
  GROUP BY bp.flight_id
)
SELECT
  (SELECT COUNT(*) FROM flights) AS total_flights,
  (SELECT COALESCE(SUM(pax)::bigint, 0) FROM pax_per_flight) AS total_passengers,
  COALESCE(AVG(p.pax / NULLIF(s.seats_per_plane, 0) * 100), 0) AS avg_load_factor
FROM seats s
LEFT JOIN pax_per_flight p ON TRUE;
`

	var stats models.AirplaneStats
	stats.AirplaneCode = code

	err := r.pool.QueryRow(ctx, query, code).Scan(
		&stats.TotalFlights,
		&stats.TotalPassengers,
		&stats.AvgLoadFactor,
	)
	if err != nil {
		return nil, fmt.Errorf("get airplane stats: %w", err)
	}
	return &stats, nil
}
