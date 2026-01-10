package handlers

import (
	"net/http"
	"time"
)

func (h *Handler) TopRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := qInt(r, "limit", 10)

	from, err := qDateRequired(r, "from")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "validation failed",
			"details": map[string]any{
				"from": err.Error(),
			},
			"example": "/api/v1/stats/routes?limit=10&from=2025-12-12&to=2026-01-11",
		})
		return
	}

	to, err := qDateRequired(r, "to")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "validation failed",
			"details": map[string]any{
				"to": err.Error(),
			},
			"example": "/api/v1/stats/routes?limit=10&from=2025-12-12&to=2026-01-11",
		})
		return
	}

	// важно: "to" делаем эксклюзивной границей конца дня
	to = to.Add(24 * time.Hour)

	stats, err := h.statsService.TopRoutes(ctx, from, to, limit)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
