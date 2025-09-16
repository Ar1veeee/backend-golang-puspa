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

type AdminHandler struct {
	adminUC usecases.AdminUseCase
}

func NewAdminHandler(adminUC usecases.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUC: adminUC,
	}
}

func (h AdminHandler) CreateAdmin(c *gin.Context) {
	req := dto.AdminCreateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	if err := h.adminUC.CreateAdminUseCase(c.Request.Context(), &req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, types.SuccessResponse{
		Success: true,
		Message: "Admin created successfully",
		Data:    nil,
	})
}

func (h AdminHandler) FindAdmins(c *gin.Context) {
	admins, err := h.adminUC.FindAdminsUseCase(c.Request.Context())
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "List of admins",
		Data:    admins,
	})
}

func (h AdminHandler) FindAdminDetail(c *gin.Context) {
	adminId := c.Param("admin_id")

	adminDetail, err := h.adminUC.FindAdminDetailUseCase(c.Request.Context(), adminId)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Admin Detail",
		Data:    adminDetail,
	})
}

func (h AdminHandler) UpdateAdmin(c *gin.Context) {
	adminId := c.Param("admin_id")

	if adminId == "" {
		middlewares.AbortWithError(c, fmt.Errorf("adminId is empty"))
		return
	}

	req := dto.AdminUpdateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	if err := h.adminUC.UpdateAdminUseCase(c.Request.Context(), adminId, &req); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Admin updated successfully",
		Data:    nil,
	})
}

func (h AdminHandler) DeleteAdmin(c *gin.Context) {
	adminId := c.Param("admin_id")

	if adminId == "" {
		middlewares.AbortWithError(c, fmt.Errorf("adminId is empty"))
		return
	}

	if err := h.adminUC.DeleteAdminWithTx(c.Request.Context(), adminId); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "Admin deleted successfully",
		Data:    nil,
	})
}
