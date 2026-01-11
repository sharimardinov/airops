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
	WITH top AS (
  SELECT f.route_no, count(*)::bigint AS flights_cnt
  FROM bookings.flights f
  WHERE f.scheduled_departure >= $1 AND f.scheduled_departure < $2
  GROUP BY f.route_no
  ORDER BY flights_cnt DESC
  LIMIT $3
)
SELECT
  r.departure_airport,
  r.arrival_airport,
  top.flights_cnt
FROM top
JOIN LATERAL (
  SELECT departure_airport, arrival_airport
  FROM bookings.routes r
  WHERE r.route_no = top.route_no
    AND r.validity && tstzrange($1, $2, '[)')   -- пересекается с периодом
  ORDER BY upper(r.validity) DESC
  LIMIT 1
) r ON true
ORDER BY top.flights_cnt DESC;
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
