package helpers

import (
	userErrors "backend-golang/internal/user/errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func IsValidPassword(password string) error {
	if len(password) < 8 {
		return userErrors.ErrWeakPassword
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return userErrors.ErrPasswordUpper
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return userErrors.ErrPasswordNumber
	}

	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	if !hasSpecial {
		return userErrors.ErrPasswordSpecial
	}

	return nil
}
