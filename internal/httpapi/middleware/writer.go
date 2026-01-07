package middleware

import "net/http"

type trackedWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (tw *trackedWriter) WriteHeader(statusCode int) {
	if tw.wroteHeader {
		return
	}
	tw.wroteHeader = true
	tw.ResponseWriter.WriteHeader(statusCode)
}

func (tw *trackedWriter) Write(b []byte) (int, error) {
	if !tw.wroteHeader {
		tw.WriteHeader(http.StatusOK)
	}
	return tw.ResponseWriter.Write(b)
}
