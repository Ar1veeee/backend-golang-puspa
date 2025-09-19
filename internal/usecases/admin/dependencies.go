package admin

import (
	"backend-golang/internal/domain/repositories"
)

type Dependencies struct {
	TxRepo    repositories.TransactionRepository
	UserRepo  repositories.UserRepository
	AdminRepo repositories.AdminRepository
	Mapper    Mapper
	Validator Validator
}

func NewDependencies(
	txRepo repositories.TransactionRepository,
	userRepo repositories.UserRepository,
	adminRepo repositories.AdminRepository,
) *Dependencies {
	return &Dependencies{
		TxRepo:    txRepo,
		UserRepo:  userRepo,
		AdminRepo: adminRepo,
		Mapper:    NewAdminMapper(),
		Validator: NewAdminValidator(),
	}
}
