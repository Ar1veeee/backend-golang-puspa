package handlers

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/adapters/http/types"
	"backend-golang/internal/usecases"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TherapistHandler struct {
	therapistUC usecases.TherapistUseCase
}

func NewTherapistHandler(therapistUC usecases.TherapistUseCase) *TherapistHandler {
	return &TherapistHandler{
		therapistUC: therapistUC,
	}
}

func (h *TherapistHandler) CreateTherapist(c *gin.Context) {
	req := dto.TherapistCreateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	err := h.therapistUC.CreateTherapistUseCase(c.Request.Context(), &req)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "Therapist Created Successfully",
		Data:    nil,
	})
}

func (h *TherapistHandler) FindTherapists(c *gin.Context) {
	therapists, err := h.therapistUC.FindTherapistsUseCase(c.Request.Context())
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "List of therapists",
		Data:    therapists,
	})
}

func (h TherapistHandler) FindTherapistDetail(c *gin.Context) {
	therapistId := c.Param("therapist_id")

	therapistDetail, err := h.therapistUC.FindTherapistDetailUseCase(c.Request.Context(), therapistId)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Therapist Detail",
		Data:    therapistDetail,
	})
}

func (h TherapistHandler) UpdateTherapist(c *gin.Context) {
	therapistId := c.Param("therapist_id")

	if therapistId == "" {
		middlewares.AbortWithError(c, fmt.Errorf("therapistId is empty"))
	}

	req := dto.TherapistUpdateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	if err := h.therapistUC.UpdateTherapistUseCase(c.Request.Context(), therapistId, &req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Therapist updated successfully",
		Data:    nil,
	})
}

func (h TherapistHandler) DeleteTherapist(c *gin.Context) {
	therapistId := c.Param("therapist_id")

	if therapistId == "" {
		middlewares.AbortWithError(c, fmt.Errorf("therapistId is empty"))
	}

	if err := h.therapistUC.DeleteTherapistWithTx(c.Request.Context(), therapistId); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Therapist deleted successfully",
		Data:    nil,
	})
}
