package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsStore struct {
	pool *pgxpool.Pool
}

func NewStatsStore(pool *pgxpool.Pool) *StatsStore {
	return &StatsStore{pool: pool}
}

type RouteStat struct {
	RouteNo          string   `json:"route_no"`
	DepartureAirport string   `json:"departure_airport"`
	ArrivalAirport   string   `json:"arrival_airport"`
	FlightsCount     int64    `json:"flights_count"`
	BoardedTotal     int64    `json:"boarded_total"`
	AvgLoadFactor    *float64 `json:"avg_load_factor,omitempty"`
	AvgDelayMinutes  *float64 `json:"avg_delay_minutes,omitempty"`
}

type RoutesStatsFilter struct {
	From     *string
	To       *string
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
	Sort     string // "boarded" | "load"
}

func (s *StatsStore) Routes(ctx context.Context, f RoutesStatsFilter) ([]RouteStat, error) {
	const qSortBoarded = `
with base as (
  select
    f.flight_id,
    f.route_no,
    r.departure_airport,
    r.arrival_airport,
    r.airplane_code,
    f.scheduled_departure,
    f.actual_departure
  from bookings.flights f
  join bookings.routes r
    on r.route_no = f.route_no
   and (f.scheduled_departure <@ r.validity)
  where ($1::text is null or r.departure_airport = $1)
    and ($2::text is null or r.arrival_airport = $2)
    and ($3::timestamptz is null or f.scheduled_departure >= $3)
    and ($4::timestamptz is null or f.scheduled_departure <  $4)
),
seats_per_aircraft as (
  select
    airplane_code,
    count(*)::float as seats_total
  from bookings.seats
  group by airplane_code
),
boarded_per_flight as (
  select
    bp.flight_id,
    count(*)::float as boarded
  from bookings.boarding_passes bp
  join base b
    on b.flight_id = bp.flight_id
  group by bp.flight_id
)
select
  b.route_no,
  b.departure_airport,
  b.arrival_airport,
  count(*)::bigint as flights_count,
  coalesce(sum(coalesce(bpf.boarded, 0)), 0)::bigint as boarded_total,
  avg(
    coalesce(bpf.boarded, 0) / sp.seats_total
  ) as avg_load_factor,
  avg(
    extract(epoch from (b.actual_departure - b.scheduled_departure)) / 60.0
  ) filter (where b.actual_departure is not null) as avg_delay_minutes
from base b
left join boarded_per_flight bpf
  on bpf.flight_id = b.flight_id
join seats_per_aircraft sp
  on sp.airplane_code = b.airplane_code
group by
  b.route_no,
  b.departure_airport,
  b.arrival_airport
order by boarded_total desc
limit $5 offset $6;
`

	const qSortLoad = `
with base as (
  select
    f.flight_id,
    f.route_no,
    r.departure_airport,
    r.arrival_airport,
    r.airplane_code,
    f.scheduled_departure,
    f.actual_departure
  from bookings.flights f
  join bookings.routes r
    on r.route_no = f.route_no
   and (f.scheduled_departure <@ r.validity)
  where ($1::text is null or r.departure_airport = $1)
    and ($2::text is null or r.arrival_airport = $2)
    and ($3::timestamptz is null or f.scheduled_departure >= $3)
    and ($4::timestamptz is null or f.scheduled_departure <  $4)
),
seats_per_aircraft as (
  select
    airplane_code,
    count(*)::float as seats_total
  from bookings.seats
  group by airplane_code
),
boarded_per_flight as (
  select
    bp.flight_id,
    count(*)::float as boarded
  from bookings.boarding_passes bp
  join base b
    on b.flight_id = bp.flight_id
  group by bp.flight_id
)
select
  b.route_no,
  b.departure_airport,
  b.arrival_airport,
  count(*)::bigint as flights_count,
  coalesce(sum(coalesce(bpf.boarded, 0)), 0)::bigint as boarded_total,
  avg(
    coalesce(bpf.boarded, 0) / sp.seats_total
  ) as avg_load_factor,
  avg(
    extract(epoch from (b.actual_departure - b.scheduled_departure)) / 60.0
  ) filter (where b.actual_departure is not null) as avg_delay_minutes
from base b
left join boarded_per_flight bpf
  on bpf.flight_id = b.flight_id
join seats_per_aircraft sp
  on sp.airplane_code = b.airplane_code
group by
  b.route_no,
  b.departure_airport,
  b.arrival_airport
order by avg_load_factor asc nulls last
limit $5 offset $6;
`

	q := qSortBoarded
	if f.Sort == "load" {
		q = qSortLoad
	}

	rows, err := s.pool.Query(ctx, q, f.From, f.To, f.DateFrom, f.DateTo, f.Limit, f.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]RouteStat, 0, f.Limit)
	for rows.Next() {
		var r RouteStat
		if err := rows.Scan(
			&r.RouteNo,
			&r.DepartureAirport,
			&r.ArrivalAirport,
			&r.FlightsCount,
			&r.BoardedTotal,
			&r.AvgLoadFactor,
			&r.AvgDelayMinutes,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
