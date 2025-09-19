package auth

import (
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/domain/services"
)

type Dependencies struct {
	TxRepo           repositories.TransactionRepository
	UserRepo         repositories.UserRepository
	ParentRepo       repositories.ParentRepository
	VerifyTokenRepo  repositories.VerificationTokenRepository
	RefreshTokenRepo repositories.RefreshTokenRepository
	EmailService     services.EmailService
	RateLimiter      services.RateLimiterService
	TokenService     services.TokenService
	Mapper           Mapper
	Validator        Validator
}

func NewDependencies(
	txRepo repositories.TransactionRepository,
	userRepo repositories.UserRepository,
	parentRepo repositories.ParentRepository,
	verifyTokenRepo repositories.VerificationTokenRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	emailService services.EmailService,
	rateLimiter services.RateLimiterService,
	tokenService services.TokenService,
) *Dependencies {
	return &Dependencies{
		TxRepo:           txRepo,
		UserRepo:         userRepo,
		ParentRepo:       parentRepo,
		VerifyTokenRepo:  verifyTokenRepo,
		RefreshTokenRepo: refreshTokenRepo,
		EmailService:     emailService,
		RateLimiter:      rateLimiter,
		TokenService:     tokenService,
		Mapper:           NewAuthMapper(),
		Validator:        NewAuthValidator(),
	}
}
