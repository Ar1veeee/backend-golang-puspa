package therapist

import "backend-golang/internal/domain/repositories"

type Dependencies struct {
	TxRepo        repositories.TransactionRepository
	UserRepo      repositories.UserRepository
	TherapistRepo repositories.TherapistRepository
	Validator     Validator
	Mapper        Mapper
}

func NewDependencies(
	txRepo repositories.TransactionRepository,
	userRepo repositories.UserRepository,
	therapistRepo repositories.TherapistRepository,
) *Dependencies {
	return &Dependencies{
		TxRepo:        txRepo,
		UserRepo:      userRepo,
		TherapistRepo: therapistRepo,
		Validator:     NewTherapistValidator(),
		Mapper:        NewTherapistMapper(),
	}
}
