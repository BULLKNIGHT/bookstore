package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/BULLKNIGHT/bookstore/logger"
)

func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(roleKey).(string)

			logger.Log.WithField("role", role).Info()

			if !ok || role != requiredRole {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
