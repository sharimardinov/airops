package handlers

import (
	"airops/internal/domain"
	"airops/internal/transport/http/middleware"
	"context"
	"errors"
	"net/http"
)

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	// error-log (всегда, но level=error)
	lg := middleware.LoggerFrom(r.Context())
	lg.Error(middleware.LogEvent{
		Msg:    "handler_error",
		RID:    middleware.GetRequestID(r.Context()),
		Method: r.Method,
		Route:  middleware.RoutePattern(r),
		Path:   r.URL.Path,
		Err:    err.Error(),
	})

	if errors.Is(err, context.Canceled) {
		// клиент сам отменил — это не ошибка сервера
		writeJSON(w, 499, apiError{Error: "client closed request"})
		return
	}

	if errors.Is(err, context.DeadlineExceeded) {
		writeJSON(w, http.StatusServiceUnavailable, apiError{Error: "timeout"})
		return
	}

	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		writeJSON(w, http.StatusBadRequest, apiError{Error: "bad request"})
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, apiError{Error: "not found"})
	default:
		writeJSON(w, http.StatusInternalServerError, apiError{Error: "internal server error"})
	}
}
