package middlewares

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/BULLKNIGHT/bookstore/logger"
	"golang.org/x/time/rate"
)

var authLimiterMap = make(map[string]*rate.Limiter)
var mutex sync.Mutex

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")

		mutex.Lock()
		limiter, exists := authLimiterMap[authToken]

		if !exists {
			limiter = rate.NewLimiter(rate.Limit(2), 10)
			authLimiterMap[authToken] = limiter
		}

		mutex.Unlock()

		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			logger.Log.Error("Too many requests")
			json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
