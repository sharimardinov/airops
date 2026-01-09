package middleware

import (
	"airops/internal/infrastructure/observability/logger"
	"net/http"
	"runtime/debug"
)

func Recover() func(http.Handler) http.Handler {
	lg := logger.NewJSONLogger()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					lg.Error(logger.LogEvent{
						Msg:    "panic",
						RID:    GetRequestID(r.Context()),
						Method: r.Method,
						Path:   r.URL.Path,
						Route:  RoutePattern(r),
						Status: 500,
						Err:    string(debug.Stack()),
					})
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
