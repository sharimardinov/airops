package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetFlightByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad flight id"})
		return
	}

	item, err := h.details.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, apiError{Error: "flight not found"})
		return
	}

	writeJSON(w, http.StatusOK, item)
}
