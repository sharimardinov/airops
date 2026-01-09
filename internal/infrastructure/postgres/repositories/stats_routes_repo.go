// stats_routes_repo.go
package repositories

import (
	"airops/internal/domain/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRoutesRepo struct {
	pool *pgxpool.Pool
}

func NewStatsRoutesRepo(pool *pgxpool.Pool) *StatsRoutesRepo {
	return &StatsRoutesRepo{pool: pool}
}

func (r *StatsRoutesRepo) TopRoutes(ctx context.Context, from, to time.Time, limit int) ([]models.RouteStat, error) {
	const q = `
select
  r.departure_airport,
  r.arrival_airport,
  count(*) as flights_cnt
from bookings.flights f
join bookings.routes r
  on r.route_no = f.route_no
 and (f.scheduled_departure <@ r.validity)
WHERE ($1::timestamptz = '0001-01-01'::timestamptz or f.scheduled_departure >= $1)
  AND ($2::timestamptz = '0001-01-01'::timestamptz or f.scheduled_departure <  $2)
group by r.departure_airport, r.arrival_airport
order by flights_cnt desc
limit $3;
`

	rows, err := r.pool.Query(ctx, q, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.RouteStat, 0, limit)

	for rows.Next() {
		var s models.RouteStat
		if err := rows.Scan(&s.DepartureAirport, &s.ArrivalAirport, &s.FlightsCount); err != nil {
			return nil, err
		}
		out = append(out, s)
	}

	return out, rows.Err()
}
