package middleware

import "net/http"

type trackedWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	bytes       int
}

func (tw *trackedWriter) WriteHeader(code int) {
	if tw.wroteHeader {
		return
	}
	tw.status = code
	tw.wroteHeader = true
	tw.ResponseWriter.WriteHeader(code)
}

func (tw *trackedWriter) Write(b []byte) (int, error) {
	if !tw.wroteHeader {
		tw.WriteHeader(http.StatusOK)
	}
	n, err := tw.ResponseWriter.Write(b)
	tw.bytes += n
	return n, err
}
