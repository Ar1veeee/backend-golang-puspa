package main

import (
	authHandler "backend-golang/internal/auth/delivery/http/handler"
	authRepo "backend-golang/internal/auth/repository/gorm"
	authService "backend-golang/internal/auth/usecase"
	registrationHandler "backend-golang/internal/registration/delivery/http/handler"
	registrationRepo "backend-golang/internal/registration/repository/gorm"
	registrationService "backend-golang/internal/registration/usecase"
	therapistHandler "backend-golang/internal/therapist/delivery/http/handler"
	therapistRepo "backend-golang/internal/therapist/repository"
	therapistService "backend-golang/internal/therapist/service"
	"backend-golang/shared/config"
	"backend-golang/shared/database"
	"backend-golang/shared/database/migrations"
	"backend-golang/shared/helpers"
	"backend-golang/shared/logger"
	"backend-golang/shared/redis"
	"backend-golang/shared/routes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnv()
	logger.InitLogger()
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	redisClient, err := redis.GetRedisClient()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	if err := helpers.InitMailjet(); err != nil {
		log.Fatalf("Failed to initialize Mailjet: %v", err)
	}

	registrationRepository := registrationRepo.NewRegistrationRepository(db)
	authRepository := authRepo.NewAuthRepository(db)
	therapistRepository := therapistRepo.NewTherapistRepository(db)

	registrationServices := registrationService.NewRegistrationService(registrationRepository)
	authServices := authService.NewAuthUseCase(authRepository, redisClient)
	therapistServices := therapistService.NewTherapistService(therapistRepository)

	registrationHandlers := registrationHandler.NewRegistrationHandler(registrationServices)
	authHandlers := authHandler.NewAuthHandler(authServices)
	therapistHandlers := therapistHandler.NewTherapistHandler(therapistServices)

	router := routes.SetupRouter(registrationHandlers, authHandlers, therapistHandlers)

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

func runMigrations(db *gorm.DB) error {
	dialect := db.Dialector.Name()
	if dialect != "mysql" {
		return fmt.Errorf("unsupported SQL dialect: %s, expected mysql", dialect)
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID:       "202509051710_create_users_table",
			Migrate:  migrations.MigrateCreateUsersTable,
			Rollback: migrations.RollbackCreateUsersTable,
		},
		{
			ID:       "202509051711_seed_admin_user",
			Migrate:  migrations.SeedUsersTableUp,
			Rollback: migrations.SeedUsersTableDown,
		},
		{
			ID:       "202509051712_create_refresh_tokens_table",
			Migrate:  migrations.MigrateCreateRefreshTokensTable,
			Rollback: migrations.RollbackCreateRefreshTokensTable,
		},
		{
			ID:       "202509051737_create_parents_table",
			Migrate:  migrations.MigrateCreateParentsTable,
			Rollback: migrations.RollbackCreateParentsTable,
		},
		{
			ID:       "202509051739_create_parent_details_table",
			Migrate:  migrations.MigrateCreateParentDetailsTable,
			Rollback: migrations.RollbackCreateParentDetailsTable,
		},
		{
			ID:       "202509051742_create_childrens_table",
			Migrate:  migrations.MigrateCreateChildrensTable,
			Rollback: migrations.RollbackCreateChildrensTable,
		},
		{
			ID:       "202509051744_create_therapists_table",
			Migrate:  migrations.MigrateCreateTherapistsTable,
			Rollback: migrations.RollbackCreateTherapistsTable,
		},
		{
			ID:       "202509071113_create_verification_codes_table",
			Migrate:  migrations.MigrateCreateVerificationCodesTable,
			Rollback: migrations.RollbackCreateVerificationCodesTable,
		},
		{
			ID:       "202509080509_create_observations_table",
			Migrate:  migrations.MigrateCreateObservationsTable,
			Rollback: migrations.RollbackCreateObservationsTable,
		},
		{
			ID:       "202509080512_create_observation_answers_table",
			Migrate:  migrations.MigrateCreateObservationAnswersTable,
			Rollback: migrations.RollbackCreateObservationAnswersTable,
		},
	})

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migrated successfully")
	return nil
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
