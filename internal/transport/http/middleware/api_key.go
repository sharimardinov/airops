package middleware

import (
	"net/http"
	"os"
)

func APIKey(key string) func(next http.Handler) http.Handler {
	if key == "" {
		key = os.Getenv("DEBUG_API_KEY")
	}
	if key == "" {
		return func(next http.Handler) http.Handler { return next }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Header.Get("X-API-Key")
			if got != key {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
