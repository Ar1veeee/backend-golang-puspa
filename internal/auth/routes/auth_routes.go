package routes

import (
	"backend-golang/internal/auth/handler"
	"backend-golang/shared/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := rg.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login",
		middlewares.RateLimiterIP(12*time.Second, 5),
		authHandler.Login,
	)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)
}
