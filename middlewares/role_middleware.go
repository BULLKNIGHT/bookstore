package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(roleKey).(string)

			fmt.Println(role)

			if !ok || role != requiredRole {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
