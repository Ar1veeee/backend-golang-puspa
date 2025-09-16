package routes

import (
	"backend-golang/internal/adapters/http/handlers"
	middlewares2 "backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/constants"
	"backend-golang/pkg/redis"
	"time"

	"github.com/gin-gonic/gin"
)

type TherapistRoutes struct {
	therapistHandler *handlers.TherapistHandler
}

func NewTherapistRoutes(
	therapistHandler *handlers.TherapistHandler,
) *TherapistRoutes {
	return &TherapistRoutes{
		therapistHandler: therapistHandler,
	}
}

func (r *TherapistRoutes) Setup(rg *gin.RouterGroup) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	therapists := rg.Group("/therapist")
	therapists.Use(
		middlewares2.Authenticate(),
		middlewares2.RateLimiterUserID(client, 1*time.Second, 100),
		middlewares2.Authorize(constants.RoleAdmin),
	)

	therapists.POST("/", r.therapistHandler.CreateTherapist)
	therapists.GET("/", r.therapistHandler.FindTherapists)
	therapists.GET("/:therapist_id", r.therapistHandler.FindTherapistDetail)
	therapists.PUT("/:therapist_id", r.therapistHandler.UpdateTherapist)
	therapists.PATCH("/:therapist_id", r.therapistHandler.DeleteTherapist)
}
