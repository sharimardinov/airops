package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"airops/internal/store"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	flights    *store.FlightsStore
	passengers *store.PassengersStore
	details    *store.FlightDetailsStore
	stats      *store.StatsStore
}

func New(pool *pgxpool.Pool) *Handler {
	return &Handler{
		flights:    store.NewFlightsStore(pool),
		passengers: store.NewPassengersStore(pool),
		details:    store.NewFlightDetailsStore(pool),
		stats:      store.NewStatsStore(pool),
	}
}

type apiError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseIntParam(r *http.Request, key string, def, min, max int) (int, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	if n < min {
		n = min
	}
	if n > max {
		n = max
	}
	return n, nil
}

func parseDateParam(r *http.Request, key string) (*time.Time, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil, nil
	}
	// ожидаем YYYY-MM-DD
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil, err
	}
	t = t.UTC()
	return &t, nil
}
