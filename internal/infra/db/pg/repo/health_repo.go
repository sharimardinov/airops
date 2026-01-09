package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthRepo struct {
	pool *pgxpool.Pool
}

func NewHealthRepo(pool *pgxpool.Pool) *HealthRepo {
	return &HealthRepo{pool: pool}
}

func (r *HealthRepo) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}
