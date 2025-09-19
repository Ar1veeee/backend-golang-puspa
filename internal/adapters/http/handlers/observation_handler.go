package handlers

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/adapters/http/types"
	"backend-golang/internal/usecases/observation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ObservationHandler struct {
	FindPendingObservationsUC   observation.FindPendingObservationsUseCase
	FindScheduledObservationsUC observation.FindScheduledObservationsUseCase
	FindCompleteObservationsUC  observation.FindCompletedObservationsUseCase
	FindObservationDetailUC     observation.FindObservationDetailUseCase
	UpdateObservationDateUC     observation.UpdateObservationDateUseCase
	ObservationQuestionsUC      observation.QuestionsUseCase
	SubmitObservationUC         observation.SubmitObservationUseCase
}

func NewObservationHandler(
	findPendingUC observation.FindPendingObservationsUseCase,
	findScheduledUC observation.FindScheduledObservationsUseCase,
	findCompleteUC observation.FindCompletedObservationsUseCase,
	findDetailUC observation.FindObservationDetailUseCase,
	updateObservationDateUC observation.UpdateObservationDateUseCase,
	observationQuestionsUC observation.QuestionsUseCase,
	submitObservationUC observation.SubmitObservationUseCase,
) *ObservationHandler {
	return &ObservationHandler{
		FindPendingObservationsUC:   findPendingUC,
		FindScheduledObservationsUC: findScheduledUC,
		FindCompleteObservationsUC:  findCompleteUC,
		FindObservationDetailUC:     findDetailUC,
		UpdateObservationDateUC:     updateObservationDateUC,
		ObservationQuestionsUC:      observationQuestionsUC,
		SubmitObservationUC:         submitObservationUC,
	}
}

func (h *ObservationHandler) FindPendingObservations(c *gin.Context) {
	observations, err := h.FindPendingObservationsUC.Execute(c.Request.Context())
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "List of pending observations",
		Data:    observations,
	})
}

func (h *ObservationHandler) FindScheduledObservations(c *gin.Context) {
	observations, err := h.FindScheduledObservationsUC.Execute(c.Request.Context())
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "List of scheduled observations",
		Data:    observations,
	})
}

func (h *ObservationHandler) FindCompletedObservations(c *gin.Context) {
	observations, err := h.FindCompleteObservationsUC.Execute(c.Request.Context())
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "List of completed observations",
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

	observationDetail, err := h.FindObservationDetailUC.Execute(c.Request.Context(), observationId)
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

func (h *ObservationHandler) UpdateObservationDate(c *gin.Context) {
	observationIdStr := c.Param("observation_id")

	observationId, err := strconv.Atoi(observationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Invalid observation ID",
		})
		return
	}

	req := dto.UpdateObservationDateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	if err := h.UpdateObservationDateUC.Execute(c.Request.Context(), observationId, &req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Observation Date Updated",
		Data:    nil,
	})
}

func (h *ObservationHandler) ObservationQuestions(c *gin.Context) {
	observationIdStr := c.Param("observation_id")

	observationId, err := strconv.Atoi(observationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Invalid observation ID",
		})
		return
	}

	observationQuestions, err := h.ObservationQuestionsUC.Execute(c.Request.Context(), observationId)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Observation Detail",
		Data:    observationQuestions,
	})
}

func (h *ObservationHandler) SubmitObservation(c *gin.Context) {
	observationIdStr := c.Param("observation_id")

	observationId, err := strconv.Atoi(observationIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Success: false,
			Message: "Invalid observation ID",
		})
		return
	}

	req := dto.SubmitObservationRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	if err := h.SubmitObservationUC.Execute(c.Request.Context(), observationId, &req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Observation Submitted",
		Data:    nil,
	})
}
