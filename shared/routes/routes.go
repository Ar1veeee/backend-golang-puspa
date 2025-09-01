package routes

import (
	authHandler "backend-golang/internal/auth/handler"
	authRoutes "backend-golang/internal/auth/routes"
	userHandler "backend-golang/internal/user/handler"
	userRoutes "backend-golang/internal/user/routes"
	"backend-golang/shared/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *authHandler.AuthHandler, userHandler *userHandler.UserHandler) *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		middlewares.RequestLogger(),
		middlewares.RateLimiterIP(10*time.Millisecond, 100),
	)

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		panic(err)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	api := router.Group("/api/v1")
	{
		authRoutes.AuthRoutes(api, authHandler)
		userRoutes.UserRoutes(api, userHandler)
	}
	return router
}
