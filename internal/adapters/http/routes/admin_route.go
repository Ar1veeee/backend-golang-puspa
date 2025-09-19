package routes

import (
	"backend-golang/internal/adapters/http/handlers"
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/constants"
	"backend-golang/pkg/redis"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminRoutes struct {
	adminHandler       *handlers.AdminHandler
	therapistHandler   *handlers.TherapistHandler
	childHandler       *handlers.ChildHandler
	observationHandler *handlers.ObservationHandler
}

func NewAdminRoutes(
	adminHandler *handlers.AdminHandler,
	therapistHandler *handlers.TherapistHandler,
	childHandler *handlers.ChildHandler,
	observationHandler *handlers.ObservationHandler,
) *AdminRoutes {
	return &AdminRoutes{
		adminHandler:       adminHandler,
		therapistHandler:   therapistHandler,
		childHandler:       childHandler,
		observationHandler: observationHandler,
	}
}

func (r *AdminRoutes) Setup(rg *gin.RouterGroup) {
	client, err := redis.GetRedisClient()
	if err != nil {
		panic(err)
	}

	admins := rg.Group("/admin")
	admins.Use(
		middlewares.Authenticate(),
		middlewares.RateLimiterUserID(client, 1*time.Second, 100),
		middlewares.Authorize(constants.RoleAdmin),
	)

	admins.POST("/admins/", r.adminHandler.CreateAdmin)
	admins.GET("/admins/", r.adminHandler.FindAdmins)
	admins.GET("/admins/:admin_id", r.adminHandler.FindAdminDetail)
	admins.PUT("/admins/:admin_id", r.adminHandler.UpdateAdmin)
	admins.PATCH("/admins/:admin_id", r.adminHandler.DeleteAdmin)

	admins.POST("/therapists/", r.therapistHandler.CreateTherapist)
	admins.GET("/therapists/", r.therapistHandler.FindTherapists)
	admins.GET("/therapists/:therapist_id", r.therapistHandler.FindTherapistDetail)
	admins.PUT("/therapists/:therapist_id", r.therapistHandler.UpdateTherapist)
	admins.PATCH("/therapists/:therapist_id", r.therapistHandler.DeleteTherapist)

	admins.GET("/childs/", r.childHandler.FindChilds)

	admins.GET("/observations/pending", r.observationHandler.FindPendingObservations)
	admins.PATCH("/observations/pending/:observation_id", r.observationHandler.UpdateObservationDate)
	admins.GET("/observations/scheduled", r.observationHandler.FindScheduledObservations)

}
