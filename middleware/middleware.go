package middleware

import (
	"github.com/lillrurre/slogr"
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter

	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func RequestLogger(logger *slogr.Logger) func(next http.Handler) http.Handler {
	logger = logger.WithGroup("request")
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.With("error", err).With("trace", string(debug.Stack())).Error("http panic")
				}
			}()

			start := time.Now()
			wrap := wrapResponseWriter(w)

			next.ServeHTTP(wrap, r)

			logger.With("duration", time.Since(start)).
				With("path", r.URL.EscapedPath()).
				With("method", r.Method).
				With("status", wrap.status).
				Info("http log")

		}
		return http.HandlerFunc(fn)
	}
}
