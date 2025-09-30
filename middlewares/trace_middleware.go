package middlewares

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the tracer from global
		tracer := otel.Tracer("bookstore-api")

		// Create a span to trace
		ctx, span := tracer.Start(r.Context(), "HTTP "+r.Method+" "+r.URL.Path)

		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.scheme", r.URL.Scheme),
			attribute.String("http.target", r.URL.Path),
			attribute.String("http.query", r.URL.RawQuery),
			attribute.String("http.user_agent", r.UserAgent()),
			attribute.String("net.peer.ip", r.RemoteAddr),
		)

		sw := &StatusWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(sw, r.WithContext(ctx))

		// End the span
		span.End()
	})
}

type StatusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (sw *StatusWriter) WriteHeader(code int) {
	sw.statusCode = code
	sw.ResponseWriter.WriteHeader(code)
}
