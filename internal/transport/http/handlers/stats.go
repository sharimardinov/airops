package handlers

import (
	"net/http"
	"time"
)

// TopRoutes GoDoc
// @Summary      Top routes
// @Description  Returns top routes for a period
// @Tags         stats
// @Accept       json
// @Produce      json
// @Param        from   query     string true "From date (YYYY-MM-DD)"
// @Param        to     query     string true "To date (YYYY-MM-DD)"
// @Param        limit  query     int    false "Limit"
// @Param from query string false "From date (YYYY-MM-DD). Default: last 30 days"
// @Param to   query string false "To date (YYYY-MM-DD). Default: today"
// @Success      200    {object}  map[string]any
// @Failure      400    {object}  map[string]any
// @Security     ApiKeyAuth
// @Router       /stats/routes [get]
func (h *Handler) TopRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit := qInt(r, "limit", 10)
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	from, fromOK, err := qDateOptional(r, "from")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error":   "validation failed",
			"details": map[string]any{"from": err.Error()},
			"example": "/api/v1/stats/routes?limit=10&from=2025-12-12&to=2026-01-11",
		})
		return
	}

	to, toOK, err := qDateOptional(r, "to")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error":   "validation failed",
			"details": map[string]any{"to": err.Error()},
			"example": "/api/v1/stats/routes?limit=10&from=2025-12-12&to=2026-01-11",
		})
		return
	}

	if !toOK {
		// сегодня 00:00 по серверному времени
		now := time.Now()
		to = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
	to = to.Add(24 * time.Hour)

	if !fromOK {
		from = to.AddDate(0, 0, -30)
	}

	stats, err := h.statsService.TopRoutes(ctx, from, to, limit)
	if err != nil {
		writeError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
