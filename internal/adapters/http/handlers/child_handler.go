package handlers

import (
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/adapters/http/types"
	"backend-golang/internal/usecases/child"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChildHandler struct {
	FindChildsUC child.FindChildUseCase
}

func NewChildHandler(
	findUC child.FindChildUseCase,

) *ChildHandler {
	return &ChildHandler{
		FindChildsUC: findUC,
	}
}

func (h ChildHandler) FindChilds(c *gin.Context) {
	childs, err := h.FindChildsUC.Execute(c.Request.Context())
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Success: true,
		Message: "List of childs",
		Data:    childs,
	})
}
