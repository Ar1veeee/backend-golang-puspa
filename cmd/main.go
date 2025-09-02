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
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.LoadEnv()
	logger.InitLogger()
	database.InitDB()

	runMigrations()

	userRepository := userRepo.NewUserRepository(database.GetDB())
	authRepository := authRepo.NewAuthRepository(database.GetDB())

	userServices := userService.NewUserService(userRepository)
	authServices := authService.NewAuthService(authRepository, userServices)

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

	setupGracefulShutdown(server)
}

func setupGracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("❌ Server forced to shutdown:", err)
	}

	if err := closeDatabase(); err != nil {
		log.Printf("❌ Error closing database: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}

func runMigrations() {
	migrationsPath := "./shared/database/migrations"

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		log.Printf("⚠️ Warning: Migrations directory not found at %s", migrationsPath)
		log.Println("⚠️ Skipping database migration")
		return
	}

	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Printf("⚠️ Warning: Cannot read migrations directory: %v", err)
		return
	}

	hasSQLFiles := false
	log.Printf("Found %d entries in migrations directory:", len(entries))
	for _, entry := range entries {
		log.Printf("- %s (isDir: %t)", entry.Name(), entry.IsDir())
		if !entry.IsDir() && len(entry.Name()) > 4 && entry.Name()[len(entry.Name())-4:] == ".sql" {
			hasSQLFiles = true
		}
	}

	if !hasSQLFiles {
		log.Println("⚠️ Warning: No .sql migration files found, skipping database migration")
		return
	}

	db, err := database.GetDB().DB()
	if err != nil {
		log.Fatal("❌ Failed to get db instance for migration: ", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal("❌ Failed to create migrate db driver: ", err)
	}

	source, err := (&file.File{}).Open("file://" + migrationsPath)
	if err != nil {
		log.Fatal("❌ Failed to create file source for migration: ", err)
	}

	m, err := migrate.NewWithInstance("file", source, "mysql", driver)
	if err != nil {
		log.Fatal("❌ Failed to create migration instance: ", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("⚠️ Warning: Could not get current migration version: %v", err)
	}

	if dirty {
		log.Printf("⚠️ Database is in dirty state at version %d. Attempting to force version...", version)
		if err := m.Force(int(version)); err != nil {
			log.Printf("❌ Failed to force migration version %d: %v", version, err)
			log.Fatal("❌ Please manually fix the database state or drop and recreate the database")
		}
		log.Printf("✅ Successfully forced database to version %d", version)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("❌ Failed to apply migrations: ", err)
	}

	log.Println("✅ Database migrated successfully")
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
