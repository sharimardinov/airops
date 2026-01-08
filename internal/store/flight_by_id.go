package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FlightDetailsStore struct {
	pool *pgxpool.Pool
}

func NewFlightDetailsStore(pool *pgxpool.Pool) *FlightDetailsStore {
	return &FlightDetailsStore{pool: pool}
}

type FlightDetails struct {
	FlightID           int64      `json:"flight_id"`
	RouteNo            string     `json:"route_no"`
	Status             string     `json:"status"`
	DepartureAirport   string     `json:"departure_airport"`
	ArrivalAirport     string     `json:"arrival_airport"`
	ScheduledDeparture time.Time  `json:"scheduled_departure"`
	ScheduledArrival   time.Time  `json:"scheduled_arrival"`
	ActualDeparture    *time.Time `json:"actual_departure,omitempty"`
	ActualArrival      *time.Time `json:"actual_arrival,omitempty"`
	BoardedCount       int64      `json:"boarded_count"`
}

func (s *FlightDetailsStore) GetByID(
	ctx context.Context,
	id int64,
) (FlightDetails, error) {
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
  f.actual_arrival,
  count(bp.ticket_no) as boarded_count
from bookings.flights f
join bookings.routes r
  on r.route_no = f.route_no
 and (f.scheduled_departure <@ r.validity)
left join bookings.boarding_passes bp
  on bp.flight_id = f.flight_id
where f.flight_id = $1
group by
  f.flight_id, f.route_no, f.status,
  r.departure_airport, r.arrival_airport,
  f.scheduled_departure, f.scheduled_arrival,
  f.actual_departure, f.actual_arrival;
`

	var d FlightDetails
	// Исправлено: использование id вместо flightID
	err := s.pool.QueryRow(ctx, q, id).Scan(
		&d.FlightID,
		&d.RouteNo,
		&d.Status,
		&d.DepartureAirport,
		&d.ArrivalAirport,
		&d.ScheduledDeparture,
		&d.ScheduledArrival,
		&d.ActualDeparture,
		&d.ActualArrival,
		&d.BoardedCount,
	)
	if err != nil {
		// Исправлено: проверка на pgx.ErrNoRows и возврат ErrNotFound
		if errors.Is(err, pgx.ErrNoRows) {
			return FlightDetails{}, ErrNotFound
		}
		return FlightDetails{}, err
	}
	// Исправлено: возврат d, а не &d
	return d, nil
}
