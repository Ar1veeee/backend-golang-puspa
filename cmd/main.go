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
	"backend-golang/shared/redis"
	"backend-golang/shared/routes"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnv()
	logger.InitLogger()
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	runMigrations(db)

	redisClient, err := redis.InitRedis()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	userRepository := userRepo.NewUserRepository(db)
	authRepository := authRepo.NewAuthRepository(db)

	userServices := userService.NewUserService(userRepository)
	authServices := authService.NewAuthService(authRepository, userServices, redisClient)

	userHandlers := userHandler.NewUserHandler(userServices)
	authHandlers := authHandler.NewAuthHandler(authServices)

	router := routes.SetupRouter(authHandlers, userHandlers)

	port := config.GetEnv("APP_PORT", "3000")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		log.Printf("🌐 Server running at http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Listen: %s\n", err)
		}
	}()

	setupGracefulShutdown(server, db)
}

func setupGracefulShutdown(server *http.Server, db *gorm.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := closeDatabase(db); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server exited gracefully")
}

func runMigrations(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get db instance for migration: ", err)
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal("Failed to create migrate db driver: ", err)
	}

	source, err := iofs.New(database.MigrationsFS, "migrations")
	if err != nil {
		log.Fatal("Failed to create iofs source for migration: ", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "mysql", driver)
	if err != nil {
		log.Fatal("Failed to create migration instance: ", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Failed to apply migrations: ", err)
	}

	log.Println("Database migrated successfully")
}

func closeDatabase(db *gorm.DB) error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
