package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Исправлено: убран status из инициализации
			tw := &trackedWriter{ResponseWriter: w}

			defer func() {
				dur := time.Since(start)
				rid := GetRequestID(r.Context())
				// status будет 200 по умолчанию, если WriteHeader не был вызван
				status := tw.status
				if status == 0 {
					status = 200
				}
				log.Printf("rid=%s %s %s status=%d bytes=%d dur=%s",
					rid, r.Method, r.URL.Path, status, tw.bytes, dur,
				)
			}()

			next.ServeHTTP(tw, r)
		})
	}
}
