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

type FlightRow struct {
	FlightID           int64     `json:"flight_id"`
	RouteNo            string    `json:"route_no"`
	DepartureAirport   string    `json:"departure_airport"`
	ArrivalAirport     string    `json:"arrival_airport"`
	ScheduledDeparture time.Time `json:"scheduled_departure"`
	Status             string    `json:"status"`
	BoardedCount       int64     `json:"boarded_count"`
}

type FlightsFilter struct {
	From     *string
	To       *string
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}

func (s *FlightsStore) List(ctx context.Context, f FlightsFilter) ([]FlightRow, error) {
	const q = `
with base as (
  select
    f.flight_id,
    f.route_no,
    r.departure_airport,
    r.arrival_airport,
    f.scheduled_departure,
    f.status
  from bookings.flights f
join bookings.routes r
  on r.route_no = f.route_no
and (f.scheduled_departure <@ r.validity)
where ($1::text is null or r.departure_airport = $1)
    and ($2::text is null or r.arrival_airport = $2)
    and ($3::timestamp is null or f.scheduled_departure >= $3)
    and ($4::timestamp is null or f.scheduled_departure <  $4)
  order by f.scheduled_departure desc
  limit $5 offset $6
)
select
  b.flight_id,
  b.route_no,
  b.departure_airport,
  b.arrival_airport,
  b.scheduled_departure,
  b.status,
  count(bp.ticket_no) as boarded_count
from base b
left join bookings.boarding_passes bp
  on bp.flight_id = b.flight_id
group by
  b.flight_id, b.route_no, b.departure_airport, b.arrival_airport, b.scheduled_departure, b.status
order by b.scheduled_departure desc;
`

	rows, err := s.pool.Query(ctx, q, f.From, f.To, f.DateFrom, f.DateTo, f.Limit, f.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]FlightRow, 0, f.Limit)
	for rows.Next() {
		var r FlightRow
		if err := rows.Scan(
			&r.FlightID,
			&r.RouteNo,
			&r.DepartureAirport,
			&r.ArrivalAirport,
			&r.ScheduledDeparture,
			&r.Status,
			&r.BoardedCount,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
