package controllers

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BULLKNIGHT/bookstore/models"
	"github.com/golang-jwt/jwt/v5"
)

// Load and decode private key
func loadPrivateKey() *rsa.PrivateKey {
	privateKeyB64 := os.Getenv("JWT_PRIVATE_KEY_B64")

	privateKeyPEM, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		log.Fatal("Failed to decode private key:", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		log.Fatal("Failed to parse private key:", err)
	}
	return privateKey
}

func loadPublicKey() *rsa.PublicKey {
	publicKeyB64 := os.Getenv("PUBLIC_KEY_BASE64")

	publicKeyPEM, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		log.Fatal("Failed to decode public key:", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		log.Fatal("Failed to parse public key:", err)
	}
	return publicKey
}

func generateJWT(username string, role string) (string, error) {
	// Create token claims
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24-hour expiry
		"iat":      time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	// Load private key
	privateKey := loadPrivateKey()

	// Sign token with your secret
	signedToken, err := token.SignedString(privateKey)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func validateUser(r *http.Request) (models.User, error) {
	// no json data send
	if r.Body == nil {
		return models.User{}, errors.New("no data found")
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	// error during parsing json data
	if err != nil {
		return models.User{}, errors.New("invalid data")
	}

	// validate required field
	if !user.IsValid() {
		return models.User{}, errors.New("all fields (name and role) are required")
	}

	return user, nil
}

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := validateUser(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	token, err := generateJWT(user.Name, user.Role)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(token)
}
