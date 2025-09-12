package handler

import (
	"backend-golang/internal/therapist/delivery/http/dto"
	"backend-golang/internal/therapist/usecase"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TherapistHandler struct {
	therapistUseCase usecase.TherapistUseCase
}

func NewTherapistHandler(therapistUseCase usecase.TherapistUseCase) *TherapistHandler {
	return &TherapistHandler{
		therapistUseCase: therapistUseCase,
	}
}

func (h *TherapistHandler) CreateTherapist(c *gin.Context) {
	req := dto.TherapistCreateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	err := h.therapistUseCase.CreateTherapistUseCase(c.Request.Context(), &req)
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

func (h *TherapistHandler) FindAllTherapists(c *gin.Context) {
	therapists, err := h.therapistUseCase.GetAllTherapistUseCase(c.Request.Context())
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

// func (h *TherapistHandler) FindUserById(c *gin.Context) {
// 	id := c.Param("id")
// 	user, err := h.therapistUseCase.GetUserByID(c.Request.Context(), id)
// 	if err != nil {
// 		middlewares.AbortWithError(c, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, types.SuccessResponse{
// 		Success: true,
// 		Message: "User Found",
// 		Data:    user,
// 	})
// }
