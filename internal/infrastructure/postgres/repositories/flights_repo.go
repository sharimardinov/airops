// flights_repo.go
package repositories

import (
	"airops/internal/domain"
	"airops/internal/domain/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FlightsRepo struct {
	pool *pgxpool.Pool
}

func NewFlightsRepo(pool *pgxpool.Pool) *FlightsRepo {
	return &FlightsRepo{pool: pool}
}

func (r *FlightsRepo) List(ctx context.Context, from, to time.Time, limit, offset int) ([]models.Flight, error) {
	const q = `
select
  f.flight_id,
  f.route_no,
  f.status,
  r.departure_airport,
  r.arrival_airport,
  f.scheduled_departure,
  f.scheduled_arrival,
  f.actual_departure,
  f.actual_arrival
from bookings.flights f
join bookings.routes r
  on r.route_no = f.route_no
 and (f.scheduled_departure <@ r.validity)
where ($1::timestamptz = '0001-01-01'::timestamptz or f.scheduled_departure >= $1)
  and ($2::timestamptz = '0001-01-01'::timestamptz or f.scheduled_departure < $2)
order by f.scheduled_departure
limit $3 offset $4;
`

	rows, err := r.pool.Query(ctx, q, from, to, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Flight, 0, limit)

	for rows.Next() {
		var f models.Flight
		if err := rows.Scan(
			&f.ID,
			&f.RouteNo,
			&f.Status,
			&f.DepartureAirport,
			&f.ArrivalAirport,
			&f.ScheduledDeparture,
			&f.ScheduledArrival,
			&f.ActualDeparture,
			&f.ActualArrival,
		); err != nil {
			return nil, err
		}
		out = append(out, f)
	}

	return out, rows.Err()
}

func (r *FlightsRepo) GetByID(ctx context.Context, id int64) (models.Flight, error) {
	const q = `
select
  f.flight_id,
  f.route_no,
  f.status,
  r.departure_airport,
  r.arrival_airport,
  f.scheduled_departure,
  f.scheduled_arrival,
  f.actual_departure,
  f.actual_arrival
from bookings.flights f
join bookings.routes r
  on r.route_no = f.route_no
 and (f.scheduled_departure <@ r.validity)
where f.flight_id = $1;
`

	var f models.Flight
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&f.ID,
		&f.RouteNo,
		&f.Status,
		&f.DepartureAirport,
		&f.ArrivalAirport,
		&f.ScheduledDeparture,
		&f.ScheduledArrival,
		&f.ActualDeparture,
		&f.ActualArrival,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Flight{}, domain.ErrNotFound
		}
		return models.Flight{}, err
	}

	return f, nil
}
