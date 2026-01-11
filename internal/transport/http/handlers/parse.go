package handlers

import (
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

func qDateOptional(r *http.Request, key string) (time.Time, bool, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return time.Time{}, false, nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return time.Time{}, true, err
	}
	return t, true, nil
}
