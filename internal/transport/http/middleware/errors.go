package middleware

import (
	"airops/internal/infrastructure/observability/logger"
	"net/http"
)

func ErrorLogging() func(http.Handler) http.Handler {
	lg := logger.NewJSONLogger()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tw := &trackedWriter{ResponseWriter: w}
			next.ServeHTTP(tw, r)

			status := tw.status
			if status == 0 {
				status = http.StatusOK
			}
			if status < 500 {
				return
			}

			rid := GetRequestID(r.Context())
			lg.Error(logger.LogEvent{
				Msg:    "server_error",
				RID:    rid,
				Method: r.Method,
				Path:   r.URL.Path,
				Status: status,
			})
		})
	}
}
