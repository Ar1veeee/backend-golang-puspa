package routes

import (
	"backend-golang/internal/adapters/http/handlers"
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/constants"
	"backend-golang/pkg/redis"
	"time"

	"github.com/gin-gonic/gin"
)

type TherapistRoutes struct {
	observationHandler *handlers.ObservationHandler
}

func NewTherapistRoutes(
	observationHandler *handlers.ObservationHandler,
) *TherapistRoutes {
	return &TherapistRoutes{
		observationHandler: observationHandler,
	}
}

func (r *TherapistRoutes) Setup(rg *gin.RouterGroup) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	therapists := rg.Group("/therapist")
	therapists.Use(
		middlewares.Authenticate(),
		middlewares.RateLimiterUserID(client, 1*time.Second, 100),
		middlewares.Authorize(constants.RoleTherapist),
	)

	therapists.GET("/observations/scheduled", r.observationHandler.FindScheduledObservations)
	therapists.GET("/observations/scheduled/:observation_id", r.observationHandler.FindObservationDetail)

	therapists.GET("/observations/question/:observation_id", r.observationHandler.ObservationQuestions)
	therapists.GET("/observations/submit/:observation_id", r.observationHandler.SubmitObservation)

	therapists.GET("/observations/completed", r.observationHandler.FindCompletedObservations)
	therapists.GET("/observations/completed/:observation_id", r.observationHandler.FindObservationDetail)
}
