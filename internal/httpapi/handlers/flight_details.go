package handlers

import (
	"airops/internal/store"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetFlightByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad id"})
		return
	}

	item, err := h.flights.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrInvalidArgument):
			writeJSON(w, http.StatusBadRequest, apiError{Error: "bad id"})
		case errors.Is(err, store.ErrNotFound):
			writeJSON(w, http.StatusNotFound, apiError{Error: "not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, apiError{Error: "internal server error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, item)
}
