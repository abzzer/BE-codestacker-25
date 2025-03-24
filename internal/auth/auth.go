package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "thisissosecretprobably"
	}
	return secret
}

func GenerateJWT(userID, role, clearance string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":         userID,
		"role":            role,
		"clearance_level": clearance,
		"exp":             time.Now().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
