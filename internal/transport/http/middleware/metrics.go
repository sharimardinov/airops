package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

func Metrics() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			tw := &trackedWriter{ResponseWriter: w}
			next.ServeHTTP(tw, r)

			status := tw.status
			if status == 0 {
				status = http.StatusOK
			}

			path := RoutePattern(r)
			for len(path) > 1 && path[0] == '/' && path[1] == '/' {
				path = path[1:]
			}

			labels := prometheus.Labels{
				"method": r.Method,
				"path":   path,
				"status": strconv.Itoa(status),
			}

			httpRequestsTotal.With(labels).Inc()
			httpRequestDuration.With(labels).Observe(time.Since(start).Seconds())
		})
	}
}
