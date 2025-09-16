package main

import (
	"backend-golang/internal/infrastructure/config"
	"backend-golang/internal/infrastructure/container"
	"backend-golang/internal/infrastructure/database"
	"backend-golang/internal/infrastructure/server"
	"backend-golang/pkg/logger"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.LoadEnv()
	logger.InitLogger()

	appContainer, err := initializeContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	defer func() {
		if err := appContainer.Close(); err != nil {
			log.Fatal("error closing container:", err)
		}
	}()

	if err := runMigrations(appContainer.DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	srv := server.NewServer(appContainer)
	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		log.Printf("🌐 Server running at http://localhost:%s", port)

		if err := srv.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Listen: %s\n", err)
		}
	}()

	setupGracefulShutdown(srv, appContainer)
}

func initializeContainer() (*container.Container, error) {
	return container.NewContainer()
}

func runMigrations(dbConn database.Connection) error {
	migrator := database.NewMigrator(dbConn.GetDB())
	return migrator.RunMigrations()
}

func setupGracefulShutdown(srv *server.Server, appContainer *container.Container) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := appContainer.Close(); err != nil {
		log.Printf("failed to close container: %v", err)
	}

	log.Println("Server exited gracefully")
}
