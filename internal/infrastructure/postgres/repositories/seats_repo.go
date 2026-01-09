package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatsRepo struct {
	pool *pgxpool.Pool
}

func NewSeatsRepo(pool *pgxpool.Pool) *SeatsRepo {
	return &SeatsRepo{pool: pool}
}

func (r *SeatsRepo) IsSeatAvailable(ctx context.Context, flightID int64, seatNo string) (bool, error) {
	seatNo = strings.TrimSpace(seatNo)
	if seatNo == "" {
		return false, fmt.Errorf("seatNo is empty")
	}

	const q = `
select not exists (
  select 1
  from bookings.boarding_passes bp
  where bp.flight_id = $1 and bp.seat_no = $2
);
`
	var ok bool
	if err := r.pool.QueryRow(ctx, q, flightID, seatNo).Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (r *SeatsRepo) GetAvailableCount(ctx context.Context, flightID int64, fareClass string) (int, error) {
	// Вариант без f.aircraft_code:
	// берём самолёт из routes (через route_no + validity)
	const q = `
select count(*)::int
from bookings.flights f
join bookings.routes r
  on r.route_no = f.route_no
 and (f.scheduled_departure <@ r.validity)
join bookings.seats s
  on s.airplane_code = r.airplane_code
where f.flight_id = $1
  and s.fare_conditions = $2
  and not exists (
    select 1
    from bookings.boarding_passes bp
    where bp.flight_id = f.flight_id
      and bp.seat_no = s.seat_no
  );
`
	var n int
	if err := r.pool.QueryRow(ctx, q, flightID, fareClass).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func (r *SeatsRepo) Reserve(ctx context.Context, tx pgx.Tx, flightID int64, ticketNo, seatNo string) error {
	seatNo = strings.TrimSpace(seatNo)
	if seatNo == "" {
		return fmt.Errorf("seatNo is empty")
	}

	const q = `
-- 1) Лочим конкретный flight_id на время транзакции (дешево и эффективно)
select pg_advisory_xact_lock($1);

-- 2) Пытаемся занять место. Если место уже занято — ничего не вставится.
with next_no as (
  select coalesce(max(boarding_no), 0) + 1 as n
  from bookings.boarding_passes
  where flight_id = $1
)
insert into bookings.boarding_passes (ticket_no, flight_id, seat_no, boarding_no, boarding_time)
select $2, $1, $3, next_no.n, now()
from next_no
on conflict (flight_id, seat_no) do nothing;
`

	ct, err := tx.Exec(ctx, q, flightID, ticketNo, seatNo)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		// место уже заняли в параллельном запросе
		return fmt.Errorf("seat %s is not available", seatNo)
	}
	return nil
}
