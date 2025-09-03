package routes

import (
	"backend-golang/internal/auth/handler"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler) {
	client, err := redis.InitRedis()
	if err != nil {
		panic(err)
	}

	auth := rg.Group("/auth")
	auth.Use(middlewares.RateLimiterIP(client, 1*time.Minute, 10))

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)
}
