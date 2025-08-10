package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func NewHTTPMw(l *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		l = l.With(
			slog.String("component", "http/middleware"),
		)
		l.Info("starting http middleware")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := l.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("request_id", middleware.GetReqID(r.Context())),
				// TODO: more info?
			)

			// Response info wrapper
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			startRequest := time.Now()
			defer func() {
				entry.Info("request completed. http middleware finished",
					slog.Int("status", ww.Status()),
					slog.Int("bytes written ", ww.BytesWritten()),
					slog.String("duration", time.Since(startRequest).String()),
				)
			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
