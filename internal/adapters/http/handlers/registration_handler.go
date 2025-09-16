package handlers

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/adapters/http/types"
	"backend-golang/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	registrationUC usecases.RegistrationUseCase
}

func NewRegistrationHandler(registrationUC usecases.RegistrationUseCase) *RegistrationHandler {
	return &RegistrationHandler{
		registrationUC: registrationUC,
	}
}

func (h *RegistrationHandler) Registration(c *gin.Context) {
	req := dto.RegistrationRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	err := h.registrationUC.RegistrationUseCase(c.Request.Context(), &req)
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
