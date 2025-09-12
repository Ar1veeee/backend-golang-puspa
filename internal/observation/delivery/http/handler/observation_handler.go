package handler

import (
	"backend-golang/internal/observation/usecase"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ObservationHandler struct {
	observationUseCase usecase.ObservationUseCase
}

func NewObservationHandler(observationUseCase usecase.ObservationUseCase) *ObservationHandler {
	return &ObservationHandler{
		observationUseCase: observationUseCase,
	}
}

func (h *ObservationHandler) FindAllObservations(c *gin.Context) {
	observations, err := h.observationUseCase.GetAllObservationsUseCase(c.Request.Context())
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "List of observations",
		Data:    observations,
	})
}
