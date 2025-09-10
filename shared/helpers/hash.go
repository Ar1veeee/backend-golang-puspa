package helpers

import (
	globalErrors "backend-golang/shared/errors"
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
		return globalErrors.ErrWeakPassword
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return globalErrors.ErrPasswordUpper
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return globalErrors.ErrPasswordNumber
	}

	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	if !hasSpecial {
		return globalErrors.ErrPasswordSpecial
	}

	return nil
}
