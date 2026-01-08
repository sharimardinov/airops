package http

import (
	"net/http"
	"time"

	"airops/internal/transport/http/handlers"
	middleware "airops/internal/transport/http/middleware"

	"github.com/go-chi/chi/v5"
)

func New(h *handlers.Handler) http.Handler {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging())
	r.Use(middleware.Recover())
	r.Use(middleware.Timeout(4 * time.Second))

	// routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/flights", func(r chi.Router) {
		r.Get("/", h.ListFlights)
		r.Get("/{id}", h.GetFlightByID)
		r.Get("/{id}/passengers", h.ListFlightPassengers)
	})

	r.Get("/stats/routes", h.TopRoutes)

	return r
}
