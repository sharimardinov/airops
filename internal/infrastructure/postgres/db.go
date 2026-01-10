package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	// Парсим конфигурацию
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	// ✅ НАСТРОЙКА CONNECTION POOL
	// Для production нагрузки увеличиваем количество соединений
	poolConfig.MaxConns = 50                      // было ~25, увеличили до 50
	poolConfig.MinConns = 10                      // минимум активных соединений
	poolConfig.MaxConnLifetime = time.Hour        // максимальное время жизни соединения
	poolConfig.MaxConnIdleTime = 30 * time.Minute // максимальное время простоя
	poolConfig.HealthCheckPeriod = time.Minute    // проверка здоровья соединений

	// Создаем pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
