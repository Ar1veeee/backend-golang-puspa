package main

import (
	authHandler "backend-golang/internal/auth/handler"
	authRepo "backend-golang/internal/auth/repository"
	authService "backend-golang/internal/auth/service"
	userHandler "backend-golang/internal/user/handler"
	userRepo "backend-golang/internal/user/repository"
	userService "backend-golang/internal/user/service"
	"backend-golang/shared/config"
	"backend-golang/shared/database"
	"backend-golang/shared/logger"
	"backend-golang/shared/routes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.LoadEnv()
	database.InitDB()
	logger.InitLogger()

	runMigrations()

	userRepository := userRepo.NewUserRepository(database.GetDB())
	authRepository := authRepo.NewAuthRepository(database.GetDB())

	userServices := userService.NewUserService(userRepository)
	authServices := authService.NewAuthService(authRepository, userServices)

	userHandlers := userHandler.NewUserHandler(userServices)
	authHandlers := authHandler.NewAuthHandler(authServices)

	setupGracefulShutdown()

	router := routes.SetupRouter(authHandlers, userHandlers)

	port := config.GetEnv("APP_PORT", "3000")

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌐 Server running at http://localhost:%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Printf("❌ Failed to start server: %v", err)
	}
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("🛑 Shutting down gracefully...")

		if err := closeDatabase(); err != nil {
			log.Printf("❌ Error closing database: %v", err)
		}
		log.Println("✅ Application stopped")
		os.Exit(0)
	}()
}

func runMigrations() {
	dbUser := config.GetEnv("DB_USER", "root")
	dbPass := config.GetEnv("DB_PASS", "")
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "3306")
	dbName := config.GetEnv("DB_NAME", "")

	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	m, err := migrate.New(
		"file://shared/database/migrations",
		dsn,
	)
	if err != nil {
		log.Fatal("Failed to create migration instance: ", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Failed to apply migrations: ", err)
	}
	log.Println("Database migrated successfully")
}

func closeDatabase() error {
	db := database.GetDB()
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
