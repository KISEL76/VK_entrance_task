package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "supersecret"))

func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func GenerateJWT(userID int, login string) (string, error) {
	ttl := time.Hour
	if envTTL := os.Getenv("JWT_TTL"); envTTL != "" {
		if d, err := time.ParseDuration(envTTL); err == nil {
			ttl = d
		}
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"login":   login,
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
