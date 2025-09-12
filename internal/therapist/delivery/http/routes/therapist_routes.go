package routes

import (
	"backend-golang/internal/therapist/delivery/http/handler"
	"backend-golang/shared/constants"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func TherapistRoutes(rg *gin.RouterGroup, therapistHandler *handler.TherapistHandler) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	users := rg.Group("/therapist")
	users.Use(
		middlewares.Authenticate(),
		middlewares.RateLimiterUserID(client, 1*time.Second, 100),
		middlewares.Authorize(constants.RoleAdmin),
	)

	users.GET("/", therapistHandler.FindAllTherapists)
	users.POST("/", therapistHandler.CreateTherapist)
	// users.GET("/:id", userHandler.FindUserById)
	// users.PUT("/:id", userHandler.UpdateUser)
	// users.DELETE("/:id", userHandler.DeleteUser)
}
