package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type wrapperWriterLogger struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrapperWriterLogger) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &wrapperWriterLogger{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rec, r)

			logger.Info("serve HTTP",
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.Int("status", rec.statusCode),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
