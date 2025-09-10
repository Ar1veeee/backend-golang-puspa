package helpers

import (
	globalErrors "backend-golang/shared/errors"
	"regexp"

	"errors"
	"fmt"
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
				errorsMap[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errorsMap[field] = "Invalid email format"
			case "min":
				errorsMap[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
			case "max":
				errorsMap[field] = fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param())
			case "alphanum":
				errorsMap[field] = fmt.Sprintf("%s can only contain letters and numbers", field)
			default:
				errorsMap[field] = fmt.Sprintf("Invalid value for %s", field)
			}
		}
	} else if err != nil {
		errorsMap["errors"] = err.Error()

		if strings.Contains(strings.ToLower(err.Error()), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				errorsMap["Username"] = globalErrors.ErrUsernameExists.Error()
			}
			if strings.Contains(err.Error(), "email") {
				errorsMap["Email"] = globalErrors.ErrEmailExists.Error()
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			errorsMap["Error"] = "Record not found"
		}
	}

	return errorsMap
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
