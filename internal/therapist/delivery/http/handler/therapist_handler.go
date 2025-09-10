package handler

import (
	"backend-golang/internal/therapist/delivery/http/dto"
	"backend-golang/internal/therapist/service"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TherapistHandler struct {
	therapistService service.TherapistService
}

func NewTherapistHandler(therapistService service.TherapistService) *TherapistHandler {
	return &TherapistHandler{
		therapistService: therapistService,
	}
}

func (h *TherapistHandler) handleErrorResponse(c *gin.Context, err error) {
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
		message = "A therapist with that email already exists."
	case errors.Is(err, globalErrors.ErrUsernameExists):
		statusCode = http.StatusConflict
		message = "A therapist with that username already exists."
	case errors.Is(err, globalErrors.ErrInvalidUserID), errors.Is(err, globalErrors.ErrUserIDRequired):
		statusCode = http.StatusBadRequest
		message = "Invalid or missing therapist ID."
	case errors.Is(err, globalErrors.ErrBadRequest), errors.Is(err, globalErrors.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		message = "Bad request."
	}

	c.JSON(statusCode, types.ErrorResponse{
		Success: false,
		Message: message,
		Errors:  helpers.TranslateErrorMessage(err),
	})
}

func (h *TherapistHandler) CreateUser(c *gin.Context) {
	req := dto.TherapistCreateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	err := h.therapistService.CreateTherapist(c.Request.Context(), &req)
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "Therapist Created Successfully",
		Data:    nil,
	})
}

func (h *TherapistHandler) FindUsers(c *gin.Context) {
	therapists, err := h.therapistService.GetAllTherapist(c.Request.Context())
	if err != nil {
		h.handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "List Data Therapists",
		Data:    therapists,
	})
}

// func (h *TherapistHandler) FindUserById(c *gin.Context) {
// 	id := c.Param("id")
// 	user, err := h.therapistService.GetUserByID(c.Request.Context(), id)
// 	if err != nil {
// 		h.handleErrorResponse(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, types.SuccessResponse{
// 		Success: true,
// 		Message: "User Found",
// 		Data:    user,
// 	})
// }
