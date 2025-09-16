package validator

import (
	internalError "backend-golang/internal/errors"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func TranslateErrorMessage(err error) map[string]string {
	errorsMap := make(map[string]string)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			switch fieldError.Tag() {
			case "required":
				errorsMap["error"] = fmt.Sprintf("%s diperlukan", field)
			case "email":
				errorsMap["error"] = "Format email tidak valid"
			case "min":
				errorsMap["error"] = fmt.Sprintf("%s minimal memiliki %s karakter", field, fieldError.Param())
			case "max":
				errorsMap["error"] = fmt.Sprintf("%s maksimal memiliki %s karakter", field, fieldError.Param())
			case "alphanum":
				errorsMap["error"] = fmt.Sprintf("Format %s tidak valid", field)
			default:
				errorsMap["error"] = fmt.Sprintf("Data untuk %s tidak valid", field)
			}
		}
	} else if err != nil {
		errorsMap["errors"] = err.Error()

		if strings.Contains(strings.ToLower(err.Error()), "duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				errorsMap["error"] = "Username sudah digunakan"
			}
			if strings.Contains(err.Error(), "email") {
				errorsMap["error"] = "Email sudah digunakan"
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			errorsMap["error"] = "Data tidak ditemukan"
		}
	}

	return errorsMap
}

func IsValidPassword(password string) error {
	if len(password) < 8 {
		return internalError.ErrPasswordTooShort
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return internalError.ErrPasswordUpper
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return internalError.ErrPasswordNumber
	}

	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	if !hasSpecial {
		return internalError.ErrPasswordSpecial
	}

	return nil
}
