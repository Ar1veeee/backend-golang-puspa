package container

import (
	"backend-golang/internal/adapters/http/handlers"
	gorm "backend-golang/internal/adapters/persistence"
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/domain/services"
	"backend-golang/internal/helpers"
	"backend-golang/internal/infrastructure/database"
	"backend-golang/internal/usecases"
	pkgredis "backend-golang/pkg/redis"

	goredis "github.com/redis/go-redis/v9"
)

type Container struct {
	DB          database.Connection
	RedisClient *goredis.Client

	// Repositories
	AdminRepo        repositories.AdminRepository
	ChildRepo        repositories.ChildRepository
	ObservationRepo  repositories.ObservationRepository
	ParentDetailRepo repositories.ParentDetailRepository
	ParentRepo       repositories.ParentRepository
	RefreshTokenRepo repositories.RefreshTokenRepository
	TherapistRepo    repositories.TherapistRepository
	TxRepo           repositories.TransactionRepository
	UserRepo         repositories.UserRepository
	VerifyTokenRepo  repositories.VerificationTokenRepository

	// Services
	emailService services.EmailService
	rateLimiter  services.RateLimiterService
	tokenService services.TokenService

	// Use Cases
	AdminUC        usecases.AdminUseCase
	AuthUC         usecases.AuthUseCase
	RegistrationUC usecases.RegistrationUseCase
	ObservationUC  usecases.ObservationUseCase
	TherapistUC    usecases.TherapistUseCase

	// Handlers
	AdminHandler        *handlers.AdminHandler
	AuthHandler         *handlers.AuthHandler
	ObservationHandler  *handlers.ObservationHandler
	RegistrationHandler *handlers.RegistrationHandler
	TherapistHandler    *handlers.TherapistHandler
}

func NewContainer() (*Container, error) {
	container := &Container{}

	if err := container.initInfrastructure(); err != nil {
		return nil, err
	}

	if err := container.initRepositories(); err != nil {
		return nil, err
	}

	if err := container.initServices(); err != nil {
		return nil, err
	}

	if err := container.initUseCases(); err != nil {
		return nil, err
	}

	if err := container.initHandlers(); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) initInfrastructure() error {
	dbConfig := database.NewConfig()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		return err
	}
	c.DB = db

	redisClient, err := pkgredis.GetRedisClient()
	if err != nil {
		return err
	}
	c.RedisClient = redisClient

	if err := helpers.InitMailjet(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initRepositories() error {
	db := c.DB.GetDB()

	c.AdminRepo = gorm.NewAdminRepository(db)
	c.ChildRepo = gorm.NewChildRepository(db)
	c.ObservationRepo = gorm.NewObservationRepository(db)
	c.ParentDetailRepo = gorm.NewParentDetailRepository(db)
	c.ParentRepo = gorm.NewParentRepository(db)
	c.RefreshTokenRepo = gorm.NewRefreshTokenRepository(db)
	c.TherapistRepo = gorm.NewTherapistRepository(db)
	c.TxRepo = gorm.NewTransactionRepository(db)
	c.UserRepo = gorm.NewUserRepository(db)
	c.VerifyTokenRepo = gorm.NewVerificationTokenRepository(db)

	return nil
}

func (c *Container) initServices() error {
	c.emailService = services.NewEmailService()

	c.rateLimiter = services.NewRateLimiterService(c.RedisClient)

	c.tokenService = services.NewTokenService()

	return nil
}

func (c *Container) initUseCases() error {
	c.AdminUC = usecases.NewAdminUseCase(c.TxRepo, c.UserRepo, c.AdminRepo)
	c.AuthUC = usecases.NewAuthUseCase(c.TxRepo, c.UserRepo, c.ParentRepo, c.VerifyTokenRepo, c.RefreshTokenRepo, c.emailService, c.rateLimiter, c.tokenService)
	c.TherapistUC = usecases.NewTherapistUseCase(c.TxRepo, c.UserRepo, c.TherapistRepo)
	c.ObservationUC = usecases.NewObservationUseCase(c.ObservationRepo)
	c.RegistrationUC = usecases.NewRegistrationUseCase(c.TxRepo, c.ParentRepo, c.ParentDetailRepo, c.ChildRepo, c.ObservationRepo)

	return nil
}

func (c *Container) initHandlers() error {
	c.AdminHandler = handlers.NewAdminHandler(c.AdminUC)
	c.AuthHandler = handlers.NewAuthHandler(c.AuthUC)
	c.TherapistHandler = handlers.NewTherapistHandler(c.TherapistUC)
	c.ObservationHandler = handlers.NewObservationHandler(c.ObservationUC)
	c.RegistrationHandler = handlers.NewRegistrationHandler(c.RegistrationUC)

	return nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
