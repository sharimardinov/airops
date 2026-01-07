package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PassengersStore struct {
	pool *pgxpool.Pool
}

func NewPassengersStore(pool *pgxpool.Pool) *PassengersStore {
	return &PassengersStore{pool: pool}
}

type PassengerRow struct {
	TicketNo       string     `json:"ticket_no"`
	PassengerName  string     `json:"passenger_name"`
	FareConditions string     `json:"fare_conditions"`
	Boarded        bool       `json:"boarded"`
	SeatNo         *string    `json:"seat_no,omitempty"`
	BoardingTime   *time.Time `json:"boarding_time,omitempty"`
}

func (s *PassengersStore) ListByFlight(ctx context.Context, flightID int64, limit, offset int) ([]PassengerRow, error) {
	const q = `
select
  t.ticket_no,
  t.passenger_name,
  sgm.fare_conditions,
  (bp.ticket_no is not null) as boarded,
  bp.seat_no,
  bp.boarding_time
from bookings.segments sgm
join bookings.tickets t
  on t.ticket_no = sgm.ticket_no
left join bookings.boarding_passes bp
  on bp.ticket_no = sgm.ticket_no
 and bp.flight_id = sgm.flight_id
where sgm.flight_id = $1
order by t.ticket_no
limit $2 offset $3;
`
	rows, err := s.pool.Query(ctx, q, flightID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]PassengerRow, 0, limit)
	for rows.Next() {
		var r PassengerRow
		if err := rows.Scan(
			&r.TicketNo,
			&r.PassengerName,
			&r.FareConditions,
			&r.Boarded,
			&r.SeatNo,
			&r.BoardingTime,
		); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
