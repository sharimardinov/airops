// passengers_repo.go
package repositories

import (
	"airops/internal/domain/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PassengersRepo struct {
	pool *pgxpool.Pool
}

func NewPassengersRepo(pool *pgxpool.Pool) *PassengersRepo {
	return &PassengersRepo{pool: pool}
}

func (r *PassengersRepo) ListByFlightID(ctx context.Context, flightID int64, limit, offset int) ([]models.FlightPassenger, error) {
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

	rows, err := r.pool.Query(ctx, q, flightID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.FlightPassenger, 0, limit)

	for rows.Next() {
		var p models.FlightPassenger
		if err := rows.Scan(
			&p.TicketNo,
			&p.PassengerName,
			&p.FareConditions,
			&p.Boarded,
			&p.SeatNo,
			&p.BoardingTime,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}

	return out, rows.Err()
}
