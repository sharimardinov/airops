package middleware

import (
	"net/http"
	"time"
)

func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// message вернётся только если реально истек таймаут
		th := http.TimeoutHandler(next, d, `{"error":"request timeout"}`)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// на таймауте будет JSON
			w.Header().Set("Content-Type", "application/json")
			th.ServeHTTP(w, r)
		})
	}
}
