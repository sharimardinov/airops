package http

import (
	"airops/internal/transport/http/handlers"
	"airops/internal/transport/http/middleware"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(h *handlers.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID())
	r.Use(middleware.Logging())
	r.Use(middleware.Recover())
	r.Use(middleware.Metrics())

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	r.Get("/health", h.Health)
	r.Get("/ready", h.Ready)
	r.Handle("/metrics", promhttp.Handler())

	// Debug (protect with API key)
	r.Route("/debug", func(r chi.Router) {
		r.Use(middleware.APIKey(""))

		r.Get("/pool", h.PoolStats)

		r.Route("/pprof", func(r chi.Router) {
			r.Get("/", pprof.Index)
			r.Get("/cmdline", pprof.Cmdline)
			r.Get("/profile", pprof.Profile)
			r.Post("/symbol", pprof.Symbol)
			r.Get("/symbol", pprof.Symbol)
			r.Get("/trace", pprof.Trace)
			r.Get("/{profile}", pprof.Index)
		})
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {

		r.Use(middleware.APIKey(""))

		// ✈Flights
		r.Get("/flights", h.ListFlights)
		r.Get("/flights/{id}", h.GetFlightByID)
		r.Get("/flights/search", h.SearchFlights)
		r.Get("/flights/{id}/passengers", h.ListFlightPassengers)

		// Airports
		r.Get("/airports", h.ListAirports)
		r.Get("/airports/search", h.SearchAirports)
		r.Get("/airports/{code}", h.GetAirport)

		// Airplanes
		r.Get("/airplanes", h.ListAirplanes)
		r.Get("/airplanes/{code}", h.GetAirplane)
		r.Get("/airplanes/{code}/seats", h.GetAirplaneWithSeats)
		r.Get("/airplanes/{code}/stats", h.GetAirplaneStats)

		// pprof
		r.Get("/pprof", pprof.Index)
		r.Get("/pprof/cmdline", pprof.Cmdline)
		r.Get("/pprof/profile", pprof.Profile)
		r.Post("/pprof/symbol", pprof.Symbol)
		r.Get("/pp	rof/symbol", pprof.Symbol)
		r.Get("/pprof/trace", pprof.Trace)
		r.Get("/pprof/{profile}", pprof.Index)

		// Bookings
		r.Post("/bookings", h.CreateBooking)
		r.Get("/bookings/{bookRef}", h.GetBooking)
		r.Delete("/bookings/{bookRef}", h.CancelBooking)

		// Passengers
		r.Get("/passengers/{passengerID}/bookings", h.GetPassengerBookings)

		// Stats
		r.Get("/stats/routes", h.TopRoutes)

	})

	return r
}
