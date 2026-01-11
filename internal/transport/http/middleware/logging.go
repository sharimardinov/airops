package middleware

import (
	"airops/internal/infrastructure/observability/logger"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func Logging() func(http.Handler) http.Handler {
	lg := logger.NewJSONLogger()
	logAll := os.Getenv("LOG_HTTP_ALL") == "1"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			tw := &trackedWriter{ResponseWriter: w}

			r = r.WithContext(logger.WithLogger(r.Context(), lg))

			next.ServeHTTP(tw, r)

			if !logAll && tw.status < 400 {
				return
			}

			status := tw.status
			if status == 0 {
				status = http.StatusOK
			}

			route := RoutePattern(r)
			path := normalizePath(r.URL.Path)

			lg.Info(logger.LogEvent{
				Msg:    "request",
				RID:    GetRequestID(r.Context()),
				Method: r.Method,
				Route:  route,
				Path:   path,
				Status: status,
				Bytes:  tw.bytes,
				DurMS:  time.Since(start).Milliseconds(),
				IP:     clientIP(r),
				UA:     r.UserAgent(),
			})
		})
	}
}

func RoutePattern(r *http.Request) string {
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		if p := rctx.RoutePattern(); p != "" {
			return p
		}
	}
	return r.URL.Path
}

func normalizePath(p string) string {
	for strings.HasPrefix(p, "//") {
		p = p[1:]
	}
	if p == "" {
		return "/"
	}
	return p
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
