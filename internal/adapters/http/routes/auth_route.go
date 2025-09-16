package routes

import (
	"backend-golang/internal/adapters/http/handlers"
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/pkg/redis"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	authHandler *handlers.AuthHandler
}

func NewAuthRoutes(
	authHandler *handlers.AuthHandler,
) *AuthRoutes {
	return &AuthRoutes{
		authHandler: authHandler,
	}
}

func (r *AuthRoutes) Setup(rg *gin.RouterGroup) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	auth := rg.Group("/auth")
	auth.Use(middlewares.RateLimiterIP(client, 1*time.Minute, 10))

	auth.POST("/login", r.authHandler.Login)
	auth.POST("/logout", r.authHandler.Logout)
	auth.POST("/refresh", r.authHandler.RefreshToken)

	auth.POST("/register", r.authHandler.Register)
	auth.GET("/verify-account", r.authHandler.VerificationAccount)

	auth.POST("/forget-password", r.authHandler.ForgetPassword)
	auth.PATCH("/reset-password", r.authHandler.ResetPassword)

	auth.GET("/resend-verification-account-email", r.authHandler.ResendVerificationAccount)
	auth.GET("/resend-reset-password-email", r.authHandler.ResendForgetPassword)
}
