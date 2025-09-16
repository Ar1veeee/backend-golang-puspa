package routes

import (
	"backend-golang/internal/adapters/http/handlers"
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/pkg/redis"
	"time"

	"github.com/gin-gonic/gin"
)

type RegistrationRoutes struct {
	registrationHandler *handlers.RegistrationHandler
}

func NewRegistrationRoutes(
	registrationHandler *handlers.RegistrationHandler,
) *RegistrationRoutes {
	return &RegistrationRoutes{
		registrationHandler: registrationHandler,
	}
}

func (r *RegistrationRoutes) Setup(rg *gin.RouterGroup) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	auth := rg.Group("/")
	auth.Use(middlewares.RateLimiterIP(client, 1*time.Minute, 10))

	auth.POST("/registration", r.registrationHandler.Registration)
}
