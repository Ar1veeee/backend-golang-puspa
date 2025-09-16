package middlewares

import (
	"backend-golang/internal/adapters/http/types"
	factoryErrors "backend-golang/internal/errors"
	internalValidator "backend-golang/internal/validator"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HandleError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		c.JSON(http.StatusUnprocessableEntity, types.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  internalValidator.TranslateErrorMessage(err),
		})
		return
	}

	var httpErr factoryErrors.HTTPError
	if errors.As(err, &httpErr) {
		c.JSON(httpErr.StatusCode, types.ErrorResponse{
			Success: false,
			Message: httpErr.Message,
			Errors:  map[string]string{"error": httpErr.UserMsg},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, types.ErrorResponse{
		Success: false,
		Message: "Internal Server Error",
		Errors:  map[string]string{"error": "Terjadi kesalahan yang tidak terduga"},
	})
}

func AbortWithError(c *gin.Context, err error) {
	HandleError(c, err)
	c.Abort()
}
