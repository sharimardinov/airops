package handlers

import (
	"airops/internal/store"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListFlightPassengers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad id"})
		return
	}

	limit := qInt(r, "limit", 100)
	offset := qInt(r, "offset", 0)

	items, err := h.passengers.ListByFlightID(r.Context(), id, limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeJSON(w, http.StatusNotFound, apiError{Error: "flight not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, apiError{Error: "internal server error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, items)
}
