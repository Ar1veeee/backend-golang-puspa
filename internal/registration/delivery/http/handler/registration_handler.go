package handler

import (
	"backend-golang/internal/registration/delivery/http/dto"
	"backend-golang/internal/registration/service"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RegistrationHandler struct {
	registrationService service.RegistrationService
}

func NewRegistrationHandler(registrationService service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{
		registrationService: registrationService,
	}
}

func (h *RegistrationHandler) handleErrorResponse(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		c.JSON(http.StatusUnprocessableEntity, types.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	statusCode := http.StatusInternalServerError
	message := "An internal server errors occurred."

	switch {
	case errors.Is(err, globalErrors.ErrUserNotFound):
		statusCode = http.StatusNotFound
		message = "User not found."
	case errors.Is(err, globalErrors.ErrEmailExists):
		statusCode = http.StatusConflict
		message = "A user with that email already exists."
	case errors.Is(err, globalErrors.ErrInvalidUserID), errors.Is(err, globalErrors.ErrUserIDRequired):
		statusCode = http.StatusBadRequest
		message = "Invalid or missing user ID."
	case errors.Is(err, globalErrors.ErrBadRequest), errors.Is(err, globalErrors.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		message = "Bad request."
	}

	c.JSON(statusCode, types.ErrorResponse{
		Success: false,
		Message: message,
		Errors:  map[string]string{"errors": message},
	})
}

func (h *RegistrationHandler) Registration(c *gin.Context) {
	req := dto.RegistrationRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	err := h.registrationService.Registration(c.Request.Context(), &req)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "User created",
		Data:    nil,
	})
}
