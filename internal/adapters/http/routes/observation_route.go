package routes

import (
	"backend-golang/internal/adapters/http/handlers"
	middlewares2 "backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/constants"
	"backend-golang/pkg/redis"
	"time"

	"github.com/gin-gonic/gin"
)

type ObservationRoutes struct {
	observationHandler *handlers.ObservationHandler
}

func NewObservationRoutes(
	observationHandler *handlers.ObservationHandler,
) *ObservationRoutes {
	return &ObservationRoutes{
		observationHandler: observationHandler,
	}
}

func (r *ObservationRoutes) Setup(rg *gin.RouterGroup) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	observations := rg.Group("/observation")
	observations.Use(
		middlewares2.Authenticate(),
		middlewares2.RateLimiterUserID(client, 1*time.Second, 100),
		middlewares2.Authorize(constants.RoleTherapist, constants.RoleAdmin),
	)

	observations.GET("/", r.observationHandler.FindPendingObservations)
	observations.GET("/:observation_id", r.observationHandler.FindObservationDetail)
	observations.GET("/completed", r.observationHandler.FindCompletedObservations)
	observations.GET("/completed/:observation_id", r.observationHandler.FindObservationDetail)
}
