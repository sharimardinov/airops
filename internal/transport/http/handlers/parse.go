package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

func pathID(r *http.Request, key string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, key), 10, 64)
}

func qInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func qDateRequired(r *http.Request, key string) (time.Time, error) {
	s := r.URL.Query().Get(key)
	if s == "" {
		return time.Time{}, fmt.Errorf("missing query param: %s", key)
	}

	// ожидаем YYYY-MM-DD
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid %s (expected YYYY-MM-DD): %w", key, err)
	}

	// делаем timestamptz-диапазон: from inclusive, to exclusive
	return t, nil
}
