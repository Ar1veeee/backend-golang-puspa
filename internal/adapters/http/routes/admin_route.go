package routes

import (
	"backend-golang/internal/adapters/http/handlers"
	middlewares2 "backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/constants"
	"backend-golang/pkg/redis"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminRoutes struct {
	adminHandler *handlers.AdminHandler
}

func NewAdminRoutes(
	adminHandler *handlers.AdminHandler,
) *AdminRoutes {
	return &AdminRoutes{
		adminHandler: adminHandler,
	}
}

func (r *AdminRoutes) Setup(rg *gin.RouterGroup) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	admins := rg.Group("/admin")
	admins.Use(
		middlewares2.Authenticate(),
		middlewares2.RateLimiterUserID(client, 1*time.Second, 100),
		middlewares2.Authorize(constants.RoleAdmin),
	)

	admins.POST("/", r.adminHandler.CreateAdmin)
	admins.GET("/", r.adminHandler.FindAdmins)
	admins.GET("/:admin_id", r.adminHandler.FindAdminDetail)
	admins.PUT("/:admin_id", r.adminHandler.UpdateAdmin)
	admins.PATCH("/:admin_id", r.adminHandler.DeleteAdmin)
}
