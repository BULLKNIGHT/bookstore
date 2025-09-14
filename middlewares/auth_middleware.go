package middlewares

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/BULLKNIGHT/bookstore/logger"
	"github.com/golang-jwt/jwt/v5"
)

func loadPublicKey() *rsa.PublicKey {
	encodedKey := os.Getenv("JWT_PUBLIC_KEY_B64")

	// Decode the token from base64
	publicKeyPEM, err := base64.StdEncoding.DecodeString(encodedKey)

	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to decode jwt public key 🔑!!")
	}

	// Parse the token
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)

	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to parse jwt public key 🔑!!")
	}

	return publicKey
}

// Validate token using custom logic
func validateTokenMethod(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, errors.New("invalid token signing method")
	}

	publicKey := loadPublicKey()
	return publicKey, nil
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse token using public key
	token, err := jwt.Parse(tokenString, validateTokenMethod)

	if err != nil || !token.Valid {
		return nil, err
	}

	// Extract claims from valid token
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Verify if token provided with correct prefix
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("missing or invalid Authorization header")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validateToken(token)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Log.WithError(err).Error(err.Error())
			json.NewEncoder(w).Encode(err)
			return
		}

		ctx := context.WithValue(r.Context(), usernameKey, claims["username"])
		ctx = context.WithValue(ctx, roleKey, claims["role"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
