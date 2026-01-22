// Package http provides HTTP server implementation and request handlers.
package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Middleware defines a function to wrap http.Handler
type Middleware func(http.Handler) http.Handler

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			logger.Info("http request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.statusCode),
				slog.Duration("duration", time.Since(start)),
				slog.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						slog.Any("error", err),
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
					)

					makeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and delegates to the underlying ResponseWriter.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// ChainMiddleware chains multiple middlewares
func ChainMiddleware(middlewares ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// HealthHandler returns a simple health check handler
func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		makeJSONResponse(w, http.StatusOK, map[string]string{"status": "healthy"})
	})
}

// ReadinessHandler returns a readiness check handler
func ReadinessHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if context is cancelled (application is shutting down)
		select {
		case <-ctx.Done():
			makeJSONResponse(w, http.StatusServiceUnavailable, map[string]string{"status": "not ready"})
			return
		default:
		}

		makeJSONResponse(w, http.StatusOK, map[string]string{"status": "ready"})
	})
}
