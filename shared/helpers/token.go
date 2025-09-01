package helpers

import (
	"backend-golang/shared/config"
	"backend-golang/shared/constants"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret_key"))

type AppClaims struct {
	Role constants.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userId string, role constants.Role) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)

	claims := &AppClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "backend_golang",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken() (string, time.Time, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", time.Time{}, err
	}
	return hex.EncodeToString(bytes), expirationTime, nil
}
