package server

import (
	"backend-golang/internal/adapters/http/middlewares"
	"backend-golang/internal/adapters/http/routes"
	"backend-golang/internal/infrastructure/container"
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router    *gin.Engine
	container *container.Container
	server    *http.Server
}

func NewServer(container *container.Container) *Server {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.RequestLogger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
	}))

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		panic(err)
	}

	srv := &Server{
		router:    router,
		container: container,
	}

	srv.setupRoutes()
	return srv
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")

	adminRoutes := routes.NewAdminRoutes(s.container.AdminHandler)
	authRoutes := routes.NewAuthRoutes(s.container.AuthHandler)
	therapistRoutes := routes.NewTherapistRoutes(s.container.TherapistHandler)
	observationRoutes := routes.NewObservationRoutes(s.container.ObservationHandler)
	registrationRoutes := routes.NewRegistrationRoutes(s.container.RegistrationHandler)

	adminRoutes.Setup(api)
	authRoutes.Setup(api)
	therapistRoutes.Setup(api)
	observationRoutes.Setup(api)
	registrationRoutes.Setup(api)

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
