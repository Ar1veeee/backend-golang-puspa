package handlers

import (
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/adapters/http/types"
	"backend-golang/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ObservationHandler struct {
	observationUC usecases.ObservationUseCase
}

func NewObservationHandler(observationUC usecases.ObservationUseCase) *ObservationHandler {
	return &ObservationHandler{
		observationUC: observationUC,
	}
}

func (h *ObservationHandler) FindPendingObservations(c *gin.Context) {
	observations, err := h.observationUC.FindPendingObservationsUseCase(c.Request.Context())
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

func (h *ObservationHandler) FindCompletedObservations(c *gin.Context) {
	observations, err := h.observationUC.FindCompletedObservationsUseCase(c.Request.Context())
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

func (h *ObservationHandler) FindObservationDetail(c *gin.Context) {
	observationIdStr := c.Param("observation_id")

	observationId, err := strconv.Atoi(observationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Invalid observation ID",
		})
		return
	}

	observationDetail, err := h.observationUC.FindObservationDetailUseCase(c.Request.Context(), observationId)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Observation Detail",
		Data:    observationDetail,
	})
}
