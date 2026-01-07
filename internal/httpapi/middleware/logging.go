package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	*trackedWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.trackedWriter.WriteHeader(code)
}

func Logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			tw := &trackedWriter{ResponseWriter: w, status: 200}

			defer func() {
				dur := time.Since(start)
				rid := GetRequestID(r.Context())
				log.Printf("rid=%s %s %s status=%d bytes=%d dur=%s",
					rid, r.Method, r.URL.Path, tw.status, tw.bytes, dur,
				)
			}()

			next.ServeHTTP(tw, r)
		})
	}
}
