package handler

import (
	"backend-golang/internal/registration/delivery/http/dto"
	"backend-golang/internal/registration/usecase"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	registrationService usecase.RegistrationService
}

func NewRegistrationHandler(registrationService usecase.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{
		registrationService: registrationService,
	}
}

func (h *RegistrationHandler) Registration(c *gin.Context) {
	req := dto.RegistrationRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	err := h.registrationService.Registration(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "User created",
		Data:    nil,
	})
}
