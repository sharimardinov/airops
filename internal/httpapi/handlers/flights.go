package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"airops/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	flights    *store.FlightsStore
	passengers *store.PassengersStore
	details    *store.FlightDetailsStore
}

func New(pool *pgxpool.Pool) *Handler {
	return &Handler{
		flights:    store.NewFlightsStore(pool),
		passengers: store.NewPassengersStore(pool),
		details:    store.NewFlightDetailsStore(pool),
	}
}

type apiError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseIntParam(r *http.Request, key string, def, min, max int) (int, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	if n < min {
		n = min
	}
	if n > max {
		n = max
	}
	return n, nil
}

func parseDateParam(r *http.Request, key string) (*time.Time, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil, nil
	}
	// ожидаем YYYY-MM-DD
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil, err
	}
	// база у тебя в +00:00, поэтому делаем UTC-полночь
	t = t.UTC()
	return &t, nil
}

func (h *Handler) ListFlights(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from") // airport code
	to := r.URL.Query().Get("to")     // airport code

	var fromPtr *string
	if from != "" {
		fromPtr = &from
	}
	var toPtr *string
	if to != "" {
		toPtr = &to
	}

	dateFrom, err := parseDateParam(r, "date_from")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad date_from (use YYYY-MM-DD)"})
		return
	}
	dateTo, err := parseDateParam(r, "date_to")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad date_to (use YYYY-MM-DD)"})
		return
	}

	if dateFrom == nil && dateTo == nil {
		now := time.Now().UTC()
		df := now.AddDate(0, 0, -30)
		dt := now
		dateFrom = &df
		dateTo = &dt
	}

	limit, err := parseIntParam(r, "limit", 50, 1, 2000)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad limit"})
		return
	}
	offset, err := parseIntParam(r, "offset", 0, 0, 1_000_000_000)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad offset"})
		return
	}

	rows, err := h.flights.List(r.Context(), store.FlightsFilter{
		From:     fromPtr,
		To:       toPtr,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":  rows,
		"limit":  limit,
		"offset": offset,
	})

}

func (h *Handler) ListFlightPassengers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad flight id"})
		return
	}

	limit, err := parseIntParam(r, "limit", 50, 1, 200)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad limit"})
		return
	}
	offset, err := parseIntParam(r, "offset", 0, 0, 1_000_000_000)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad offset"})
		return
	}

	items, err := h.passengers.ListByFlight(r.Context(), id, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"flight_id": id,
		"items":     items,
		"limit":     limit,
		"offset":    offset,
	})
}
func (h *Handler) GetFlightByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad flight id"})
		return
	}

	item, err := h.details.GetByID(r.Context(), id)
	if err != nil {
		// самый простой вариант без pgx-специфики:
		writeJSON(w, http.StatusNotFound, apiError{Error: "flight not found"})
		return
	}

	writeJSON(w, http.StatusOK, item)
}
