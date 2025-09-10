package routes

import (
	authHandler "backend-golang/internal/auth/delivery/http/handler"
	authRoutes "backend-golang/internal/auth/delivery/http/routes"
	registrationHandler "backend-golang/internal/registration/delivery/http/handler"
	registrationRoutes "backend-golang/internal/registration/delivery/http/routes"
	therapistHandler "backend-golang/internal/therapist/delivery/http/handler"
	therapistRoutes "backend-golang/internal/therapist/delivery/http/routes"
	"backend-golang/shared/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	registrationHandler *registrationHandler.RegistrationHandler,
	authHandler *authHandler.AuthHandler,
	therapistHandler *therapistHandler.TherapistHandler,
) *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		middlewares.RequestLogger(),
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
		registrationRoutes.RegistrationRoutes(api, registrationHandler)
		authRoutes.AuthRoutes(api, authHandler)
		therapistRoutes.TherapistRoutes(api, therapistHandler)
	}
	return router
}
