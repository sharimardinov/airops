package httpapi

import (
	"airops/internal/httpapi/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	h := handlers.New(pool)
	r.Get("/flights/{id}", h.GetFlightByID)
	r.Get("/flights/{id}/passengers", h.ListFlightPassengers)

	r.Get("/flights", h.ListFlights)

	r.Get("/stats/routes", h.RoutesStats)

	return r
}
