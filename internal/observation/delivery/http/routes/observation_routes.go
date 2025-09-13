package routes

import (
	"backend-golang/internal/observation/delivery/http/handler"
	"backend-golang/shared/constants"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func ObservationRoutes(rg *gin.RouterGroup, observationHandler *handler.ObservationHandler) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	observations := rg.Group("/observation")
	observations.Use(
		middlewares.Authenticate(),
		middlewares.RateLimiterUserID(client, 1*time.Second, 100),
		middlewares.Authorize(constants.RoleTherapist, constants.RoleAdmin),
	)

	observations.GET("/", observationHandler.PendingObservations)
	observations.GET("/:observation_id", observationHandler.DetailObservation)
	observations.GET("/completed", observationHandler.CompletedObservations)
}
