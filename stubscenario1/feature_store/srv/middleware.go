package srv

import (
	"log"
	"net/http"
	"time"
)

// middleware
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWrite(w)
		next(rw, r)
		log.Printf("status_code:, %d, duratuion: %d ms\n", rw.status, time.Since(start).Milliseconds())
	})
}

// Hook .
func newResponseWrite(w http.ResponseWriter) *responseWriterDelegator {
	return &responseWriterDelegator{ResponseWriter: w}
}

// Hook .
type responseWriterDelegator struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

// WriteHeader Hook.
func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

// Write Hook .
func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}
