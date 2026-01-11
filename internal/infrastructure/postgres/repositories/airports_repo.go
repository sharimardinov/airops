package repositories

import (
	"airops/internal/app/apperr"
	"airops/internal/domain/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AirportsRepo struct {
	pool *pgxpool.Pool
}

func NewAirportsRepo(pool *pgxpool.Pool) *AirportsRepo {
	return &AirportsRepo{pool: pool}
}

// получает аэропорт по коду
func (r *AirportsRepo) GetByCode(ctx context.Context, code string) (*models.Airport, error) {
	query := `
		SELECT airport_code, airport_name, city, country, timezone
		FROM bookings.airports
		WHERE airport_code = $1
	`

	var airport models.Airport
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&airport.Code,
		&airport.Name,
		&airport.City,
		&airport.Country,
		&airport.Timezone,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("get airport by code", err)
	}
	if err != nil {
		return nil, apperr.Internal("db error", err)
	}

	return &airport, nil
}

// возвращает список всех аэропортов
func (r *AirportsRepo) List(ctx context.Context) ([]models.Airport, error) {
	query := `
		SELECT airport_code, airport_name, city, country, timezone
		FROM bookings.airports
		ORDER BY city, airport_name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list airports: %w", err)
	}
	defer rows.Close()

	var airports []models.Airport
	for rows.Next() {
		var airport models.Airport
		err := rows.Scan(
			&airport.Code,
			&airport.Name,
			&airport.City,
			&airport.Country,
			&airport.Timezone,
		)
		if err != nil {
			return nil, fmt.Errorf("scan airport: %w", err)
		}
		airports = append(airports, airport)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return airports, nil
}

// ищет аэропорты по названию города
func (r *AirportsRepo) SearchByCity(ctx context.Context, city string) ([]models.Airport, error) {
	query := `
		SELECT airport_code, airport_name, city, country, timezone
		FROM bookings.airports
		WHERE LOWER(city) LIKE LOWER($1)
		ORDER BY city, airport_name
	`

	rows, err := r.pool.Query(ctx, query, "%"+city+"%")
	if err != nil {
		return nil, fmt.Errorf("search airports by city: %w", err)
	}
	defer rows.Close()

	var airports []models.Airport
	for rows.Next() {
		var airport models.Airport
		err := rows.Scan(
			&airport.Code,
			&airport.Name,
			&airport.City,
			&airport.Country,
			&airport.Timezone,
		)
		if err != nil {
			return nil, fmt.Errorf("scan airport: %w", err)
		}
		airports = append(airports, airport)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return airports, nil
}

// ищет аэропорты по стране
func (r *AirportsRepo) SearchByCountry(ctx context.Context, country string) ([]models.Airport, error) {
	query := `
  SELECT airport_code, airport_name, city, country, timezone
  FROM bookings.airports
  ORDER BY city, airport_name
	`

	rows, err := r.pool.Query(ctx, query, "%"+country+"%")
	if err != nil {
		return nil, fmt.Errorf("search airports by country: %w", err)
	}
	defer rows.Close()

	var airports []models.Airport
	for rows.Next() {
		var airport models.Airport
		err := rows.Scan(
			&airport.Code,
			&airport.Name,
			&airport.City,
			&airport.Country,
			&airport.Timezone,
		)
		if err != nil {
			return nil, fmt.Errorf("scan airport: %w", err)
		}
		airports = append(airports, airport)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return airports, nil
}
