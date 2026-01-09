// internal/domain/http/router.go
package http

import (
	"airops/internal/transport/http/handlers"
	"airops/internal/transport/http/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(h *handlers.Handler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging())
	r.Use(middleware.Recover())
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Metrics())

	// Health checks
	r.Get("/health", h.Health)
	r.Get("/ready", h.Ready)
	r.Handle("/metrics", promhttp.Handler())

	r.Get("/flights", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/flights", http.StatusTemporaryRedirect)
	})
	r.Get("/airports", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/airports", http.StatusTemporaryRedirect)
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Flights
		r.Get("/flights", h.ListFlights)
		r.Get("/flights/{id}", h.GetFlightByID)
		r.Get("/flights/search", h.SearchFlights) // NEW

		// Airports
		r.Get("/airports", h.ListAirports)          // NEW
		r.Get("/airports/search", h.SearchAirports) // NEW
		r.Get("/airports/{code}", h.GetAirport)     // NEW

		// Passengers
		r.Get("/flights/{id}/passengers", h.ListFlightPassengers)

		// Stats
		r.Get("/stats/routes", h.TopRoutes)

		// Bookings
		r.Post("/bookings", h.CreateBooking)             // NEW
		r.Get("/bookings/{bookRef}", h.GetBooking)       // NEW
		r.Delete("/bookings/{bookRef}", h.CancelBooking) // NEW

		// Passenger bookings
		r.Get("/passengers/{passengerID}/bookings", h.GetPassengerBookings) // NEW
	})

	return r
}
