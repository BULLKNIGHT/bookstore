package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/BULLKNIGHT/bookstore/logger"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log := logger.Log.WithFields(map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		})

		ctx := context.WithValue(r.Context(), loggerKey, log)

		next.ServeHTTP(w, r.WithContext(ctx))

		log.WithField("duration", time.Since(start).Milliseconds()).Info("request completed")
	})
}
