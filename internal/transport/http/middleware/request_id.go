package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"
const requestIDHeader = "X-Request-Id"

func newRequestID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(requestIDHeader)
			if id == "" {
				id = newRequestID()
			}

			ctx := context.WithValue(r.Context(), requestIDKey, id)
			r = r.WithContext(ctx)
			w.Header().Set(requestIDHeader, id)
			next.ServeHTTP(w, r)
		})
	}
}

func GetRequestID(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}
