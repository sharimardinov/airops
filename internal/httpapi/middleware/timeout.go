package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			tw := &trackedWriter{ResponseWriter: w}

			// даём вниз ctx с дедлайном
			next.ServeHTTP(tw, r.WithContext(ctx))

			// если хэндлер уже ответил — не лезем
			if tw.wroteHeader {
				return
			}

			// если запрос умер по таймауту/отмене — отдаём 504
			if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.Canceled) {
				tw.Header().Set("Content-Type", "application/json")
				tw.WriteHeader(http.StatusGatewayTimeout)
				_ = json.NewEncoder(tw).Encode(map[string]string{
					"error": "request timeout",
				})
				return
			}
		})
	}
}
