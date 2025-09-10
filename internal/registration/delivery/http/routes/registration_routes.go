package routes

import (
	"backend-golang/internal/registration/delivery/http/handler"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func RegistrationRoutes(rg *gin.RouterGroup, registrationHandler *handler.RegistrationHandler) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	auth := rg.Group("/")
	auth.Use(middlewares.RateLimiterIP(client, 1*time.Minute, 10))

	auth.POST("/registration", registrationHandler.Registration)
}
