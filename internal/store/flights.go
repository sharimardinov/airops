package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FlightsStore struct {
	pool *pgxpool.Pool
}

func NewFlightsStore(pool *pgxpool.Pool) *FlightsStore {
	return &FlightsStore{pool: pool}
}

// Добавлен отсутствующий тип Flight
type Flight struct {
	FlightID           int64      `json:"flight_id"`
	RouteNo            string     `json:"route_no"`
	Status             string     `json:"status"`
	DepartureAirport   string     `json:"departure_airport"`
	ArrivalAirport     string     `json:"arrival_airport"`
	ScheduledDeparture time.Time  `json:"scheduled_departure"`
	ScheduledArrival   time.Time  `json:"scheduled_arrival"`
	ActualDeparture    *time.Time `json:"actual_departure,omitempty"`
	ActualArrival      *time.Time `json:"actual_arrival,omitempty"`
}

func (s *FlightsStore) List(
	ctx context.Context,
	from, to time.Time,
	limit, offset int,
) ([]Flight, error) {
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

	rows, err := s.pool.Query(ctx, q, from, to, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Flight, 0, limit)
	for rows.Next() {
		var flight Flight
		if err := rows.Scan(
			&flight.FlightID,
			&flight.RouteNo,
			&flight.Status,
			&flight.DepartureAirport,
			&flight.ArrivalAirport,
			&flight.ScheduledDeparture,
			&flight.ScheduledArrival,
			&flight.ActualDeparture,
			&flight.ActualArrival,
		); err != nil {
			return nil, err
		}
		out = append(out, flight)
	}
	return out, rows.Err()
}
