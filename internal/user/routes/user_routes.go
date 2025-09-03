package routes

import (
	"backend-golang/internal/user/handler"
	"backend-golang/shared/constants"
	"backend-golang/shared/middlewares"
	"backend-golang/shared/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	client, err := redis.InitRedis()
	if err != nil {
		panic(err)
	}

	users := rg.Group("/users")
	users.Use(
		middlewares.Authenticate(),
		middlewares.RateLimiterUserID(client, 1*time.Second, 100),
		middlewares.Authorize(constants.RoleAdmin),
	)

	users.GET("/", userHandler.FindUsers)
	users.POST("/", userHandler.CreateUser)
	users.GET("/:id", userHandler.FindUserById)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)
}
