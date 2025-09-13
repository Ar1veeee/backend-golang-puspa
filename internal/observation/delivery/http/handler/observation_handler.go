package handler

import (
	observationDto "backend-golang/internal/observation/delivery/http/dto"
	"backend-golang/internal/observation/usecase"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/types"
	"net/http"
	"strconv"

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

func (h *ObservationHandler) PendingObservations(c *gin.Context) {
	observations, err := h.observationUseCase.GetPendingObservationsUseCase(c.Request.Context())
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

func (h *ObservationHandler) CompletedObservations(c *gin.Context) {
	observations, err := h.observationUseCase.GetCompletedObservationsUseCase(c.Request.Context())
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

func (h *ObservationHandler) DetailObservation(c *gin.Context) {
	observationIdStr := c.Param("observation_id")

	observationId, err := strconv.Atoi(observationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Invalid observation ID",
		})
		return
	}

	req := observationDto.DetailObservationRequest{
		ObservationId: observationId,
	}

	observationDetail, err := h.observationUseCase.GetObservationDetail(c.Request.Context(), req.ObservationId)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Detail observation",
		Data:    observationDetail,
	})
}
