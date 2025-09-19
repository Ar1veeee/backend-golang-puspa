package container

import (
	"backend-golang/internal/adapters/http/handlers"
	gorm "backend-golang/internal/adapters/persistence"
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/domain/services"
	"backend-golang/internal/helpers"
	"backend-golang/internal/infrastructure/database"
	"backend-golang/internal/usecases/admin"
	"backend-golang/internal/usecases/auth"
	"backend-golang/internal/usecases/child"
	"backend-golang/internal/usecases/observation"
	"backend-golang/internal/usecases/registration"
	"backend-golang/internal/usecases/therapist"
	pkgredis "backend-golang/pkg/redis"

	goredis "github.com/redis/go-redis/v9"
)

type Container struct {
	DB          database.Connection
	RedisClient *goredis.Client

	// Repositories
	AdminRepo               repositories.AdminRepository
	ChildRepo               repositories.ChildRepository
	ObservationRepo         repositories.ObservationRepository
	ObservationQuestionRepo repositories.ObservationQuestionRepository
	ObservationAnswerRepo   repositories.ObservationAnswerRepository
	ParentDetailRepo        repositories.ParentDetailRepository
	ParentRepo              repositories.ParentRepository
	RefreshTokenRepo        repositories.RefreshTokenRepository
	TherapistRepo           repositories.TherapistRepository
	TxRepo                  repositories.TransactionRepository
	UserRepo                repositories.UserRepository
	VerifyTokenRepo         repositories.VerificationTokenRepository

	// Services
	emailService services.EmailService
	rateLimiter  services.RateLimiterService
	tokenService services.TokenService

	// Use Case Auth
	RegisterUC                  auth.RegisterUseCase
	LoginUC                     auth.LoginUseCase
	ResetPasswordUC             auth.ResetPasswordUseCase
	RefreshTokenUC              auth.RefreshTokenUseCase
	LogoutUC                    auth.LogoutUseCase
	ResendVerificationAccountUC auth.ResendVerificationAccountUseCase
	VerificationAccountUC       auth.VerificationAccountUseCase
	ForgetPasswordUC            auth.ForgetPasswordUseCase
	ResendForgetPasswordUC      auth.ResendForgetPasswordUseCase

	// Use Cases Admin
	CreateAdminUC     admin.CreateAdminUseCase
	FindAdminsUC      admin.FindAdminsUseCase
	FindAdminDetailUC admin.FindAdminDetailUseCase
	UpdateAdminUC     admin.UpdateAdminUseCase
	DeleteAdminUC     admin.DeleteAdminUseCase

	// Use Cases Therapist
	CreateTherapistUC     therapist.CreateTherapistUseCase
	FindTherapistsUC      therapist.FindTherapistsUseCase
	FindTherapistDetailUC therapist.FindTherapistDetailUseCase
	UpdateTherapistUC     therapist.UpdateTherapistUseCase
	DeleteTherapistUC     therapist.DeleteTherapistUseCase

	// Use Case Registration
	RegistrationUC registration.RegistrationUseCase

	// Use Case Child
	FindChildsUC child.FindChildUseCase

	//Use Case Observation
	FindPendingObservationsUC   observation.FindPendingObservationsUseCase
	FindScheduledObservationsUC observation.FindScheduledObservationsUseCase
	FindCompletedObservationsUC observation.FindCompletedObservationsUseCase
	FindObservationDetailUC     observation.FindObservationDetailUseCase
	UpdateObservationDateUC     observation.UpdateObservationDateUseCase
	ObservationQuestionsUC      observation.QuestionsUseCase
	SubmitObservationUC         observation.SubmitObservationUseCase

	// Handlers
	AdminHandler        *handlers.AdminHandler
	AuthHandler         *handlers.AuthHandler
	ObservationHandler  *handlers.ObservationHandler
	RegistrationHandler *handlers.RegistrationHandler
	TherapistHandler    *handlers.TherapistHandler
	ChildHandler        *handlers.ChildHandler
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
	c.ObservationQuestionRepo = gorm.NewObservationQuestionRepository(db)
	c.ObservationAnswerRepo = gorm.NewObservationAnswerRepository(db)
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
	// Auth Use Case
	authDeps := auth.NewDependencies(
		c.TxRepo,
		c.UserRepo,
		c.ParentRepo,
		c.VerifyTokenRepo,
		c.RefreshTokenRepo,
		c.emailService,
		c.rateLimiter,
		c.tokenService,
	)

	c.RegisterUC = auth.NewRegisterUseCase(authDeps)
	c.LoginUC = auth.NewLoginUseCase(authDeps)
	c.ResetPasswordUC = auth.NewResetPasswordUseCase(authDeps)
	c.RefreshTokenUC = auth.NewRefreshTokenUseCase(authDeps)
	c.LogoutUC = auth.NewLogoutUseCase(authDeps)
	c.ResendVerificationAccountUC = auth.NewResendVerificationAccountUseCase(authDeps)
	c.VerificationAccountUC = auth.NewVerificationAccountUseCase(authDeps)
	c.ForgetPasswordUC = auth.NewForgetPasswordUseCase(authDeps)
	c.ResendForgetPasswordUC = auth.NewResendForgetPasswordUseCase(authDeps)

	// Admin Use Case
	adminDeps := admin.NewDependencies(
		c.TxRepo,
		c.UserRepo,
		c.AdminRepo,
	)

	c.CreateAdminUC = admin.NewCreateAdminUseCase(adminDeps)
	c.FindAdminsUC = admin.NewFindAdminsUseCase(adminDeps)
	c.FindAdminDetailUC = admin.NewFindAdminDetailUseCase(adminDeps)
	c.UpdateAdminUC = admin.NewUpdateAdminUseCase(adminDeps)
	c.DeleteAdminUC = admin.NewDeleteAdminUseCase(adminDeps)

	// Therapist Use Case
	therapistDeps := therapist.NewDependencies(c.TxRepo, c.UserRepo, c.TherapistRepo)

	c.CreateTherapistUC = therapist.NewCreateTherapistUseCase(therapistDeps)
	c.FindTherapistsUC = therapist.NewFindTherapistsUseCase(therapistDeps)
	c.FindTherapistDetailUC = therapist.NewFindTherapistDetailUseCase(therapistDeps)
	c.UpdateTherapistUC = therapist.NewUpdateTherapistUseCase(therapistDeps)
	c.DeleteTherapistUC = therapist.NewDeleteTherapistUseCase(therapistDeps)

	// Registration Use Case
	registrationDeps := registration.NewDependencies(
		c.TxRepo,
		c.ParentRepo,
		c.ParentDetailRepo,
		c.ChildRepo,
		c.ObservationRepo,
	)

	c.RegistrationUC = registration.NewRegistrationUseCase(registrationDeps)

	// Child Use Case
	childDeps := child.NewDependencies(c.ChildRepo)

	c.FindChildsUC = child.NewFindChildUseCase(childDeps)

	// Observation Use Case
	observationDeps := observation.NewDependencies(
		c.TxRepo,
		c.ObservationRepo,
		c.ObservationQuestionRepo,
		c.ObservationAnswerRepo,
		c.TherapistRepo,
	)

	c.FindPendingObservationsUC = observation.NewFindPendingObservationsUseCase(observationDeps)
	c.FindScheduledObservationsUC = observation.NewFindScheduledObservationsUseCase(observationDeps)
	c.FindCompletedObservationsUC = observation.NewFindCompletedObservationsUseCase(observationDeps)
	c.FindObservationDetailUC = observation.NewFindObservationDetailUseCase(observationDeps)
	c.UpdateObservationDateUC = observation.NewUpdateObservationDateUseCase(observationDeps)
	c.ObservationQuestionsUC = observation.NewObservationQuestionsUseCase(observationDeps)
	c.SubmitObservationUC = observation.NewSubmitObservationUseCase(observationDeps)

	return nil
}

func (c *Container) initHandlers() error {
	c.AdminHandler = handlers.NewAdminHandler(
		c.CreateAdminUC,
		c.FindAdminsUC,
		c.FindAdminDetailUC,
		c.UpdateAdminUC,
		c.DeleteAdminUC,
	)

	c.AuthHandler = handlers.NewAuthHandler(
		c.RegisterUC,
		c.LoginUC,
		c.ResetPasswordUC,
		c.RefreshTokenUC,
		c.LogoutUC,
		c.ResendVerificationAccountUC,
		c.VerificationAccountUC,
		c.ForgetPasswordUC,
		c.ResendForgetPasswordUC,
	)

	c.TherapistHandler = handlers.NewTherapistHandler(
		c.CreateTherapistUC,
		c.FindTherapistsUC,
		c.FindTherapistDetailUC,
		c.UpdateTherapistUC,
		c.DeleteTherapistUC,
	)

	c.ObservationHandler = handlers.NewObservationHandler(
		c.FindPendingObservationsUC,
		c.FindScheduledObservationsUC,
		c.FindCompletedObservationsUC,
		c.FindObservationDetailUC,
		c.UpdateObservationDateUC,
		c.ObservationQuestionsUC,
		c.SubmitObservationUC,
	)

	c.RegistrationHandler = handlers.NewRegistrationHandler(
		c.RegistrationUC,
	)

	c.ChildHandler = handlers.NewChildHandler(
		c.FindChildsUC,
	)

	return nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
