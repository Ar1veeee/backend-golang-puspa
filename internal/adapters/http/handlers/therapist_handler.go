package handlers

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/adapters/http/types"
	"backend-golang/internal/usecases/therapist"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TherapistHandler struct {
	CreateTherapistUC     therapist.CreateTherapistUseCase
	FindTherapistsUC      therapist.FindTherapistsUseCase
	FindTherapistDetailUC therapist.FindTherapistDetailUseCase
	UpdateTherapistUC     therapist.UpdateTherapistUseCase
	DeleteTherapistUC     therapist.DeleteTherapistUseCase
}

func NewTherapistHandler(
	createUC therapist.CreateTherapistUseCase,
	findUC therapist.FindTherapistsUseCase,
	findDetailUC therapist.FindTherapistDetailUseCase,
	updateUC therapist.UpdateTherapistUseCase,
	deleteUC therapist.DeleteTherapistUseCase,
) *TherapistHandler {
	return &TherapistHandler{
		CreateTherapistUC:     createUC,
		FindTherapistsUC:      findUC,
		FindTherapistDetailUC: findDetailUC,
		UpdateTherapistUC:     updateUC,
		DeleteTherapistUC:     deleteUC,
	}
}

func (h *TherapistHandler) CreateTherapist(c *gin.Context) {
	req := dto.TherapistCreateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	err := h.CreateTherapistUC.Execute(c.Request.Context(), &req)
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
	therapists, err := h.FindTherapistsUC.Execute(c.Request.Context())
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

	therapistDetail, err := h.FindTherapistDetailUC.Execute(c.Request.Context(), therapistId)
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

	if err := h.UpdateTherapistUC.Execute(c.Request.Context(), therapistId, &req); err != nil {
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

	if err := h.DeleteTherapistUC.Execute(c.Request.Context(), therapistId); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Therapist deleted successfully",
		Data:    nil,
	})
}
