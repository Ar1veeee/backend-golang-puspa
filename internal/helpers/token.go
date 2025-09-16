package helpers

import (
	"backend-golang/internal/constants"
	"backend-golang/internal/infrastructure/config"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
	return token.SignedString(config.JWTKey)
}

func GenerateVerificationToken(userId string) (string, time.Time, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "backend_golang",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTKey)
	return tokenString, expirationTime, err
}

func VerifyVerificationToken(tokenString string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWTKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
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
