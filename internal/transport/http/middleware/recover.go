package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

func Recover() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tw := &trackedWriter{ResponseWriter: w}

			defer func() {
				if rec := recover(); rec != nil {
					// Логи — в консоль (потом заменим на нормальный логгер)
					log.Printf("panic: %v", rec)

					// Если уже начали писать ответ — поздно, просто рвёмся
					if tw.wroteHeader {
						return
					}

					tw.Header().Set("Content-Type", "application/json")
					tw.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(tw).Encode(map[string]string{
						"error": "internal server error",
					})
				}
			}()

			next.ServeHTTP(tw, r)
		})
	}
}
