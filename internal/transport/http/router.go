package http

import (
	"airops/internal/transport/http/handlers"
	middleware "airops/internal/transport/http/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(h *handlers.Handler) http.Handler {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Metrics())
	r.Use(middleware.Logging())
	r.Use(middleware.Recover())
	r.Use(middleware.ErrorLogging())

	// routes

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/health", h.Health)
	r.Get("/ready", h.Ready)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(20 * time.Second)) // или 10s

		r.Route("/flights", func(r chi.Router) {
			r.Get("/", h.ListFlights)
			r.Get("/{id}", h.GetFlightByID)
			r.Get("/{id}/passengers", h.ListFlightPassengers)
		})

		r.Get("/stats/routes", h.TopRoutes)
	})
	return r
}
