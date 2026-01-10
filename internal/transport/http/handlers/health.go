package handlers

import (
	"net/http"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.healthService == nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "health service not wired"})
		return
	}

	if err := h.healthService.Ready(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, apiError{Error: "not ready"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) PoolStats(w http.ResponseWriter, r *http.Request) {
	if h.pool == nil {
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "db pool not wired"})
		return
	}

	stats := h.pool.Stat()

	response := map[string]interface{}{
		"acquired_conns":      stats.AcquiredConns(),
		"idle_conns":          stats.IdleConns(),
		"total_conns":         stats.TotalConns(),
		"max_conns":           stats.MaxConns(),
		"new_conns_count":     stats.NewConnsCount(),
		"acquire_count":       stats.AcquireCount(),
		"acquire_duration":    stats.AcquireDuration().String(),
		"empty_acquire_count": stats.EmptyAcquireCount(),
	}

	writeJSON(w, http.StatusOK, response)
}
