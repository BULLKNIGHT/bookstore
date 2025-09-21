package middlewares

import (
	"net/http"

	"go.opentelemetry.io/otel"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the tracer from global
		tracer := otel.Tracer("bookstore-api")

		// Create a span to trace
		ctx, span := tracer.Start(r.Context(), "HTTP "+r.Method+" "+r.URL.Path)

		next.ServeHTTP(w, r.WithContext(ctx))

		// End the span
		span.End()
	})
}
