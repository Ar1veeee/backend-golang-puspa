package registration

import "backend-golang/internal/domain/repositories"

type Dependencies struct {
	TxRepo           repositories.TransactionRepository
	ParentRepo       repositories.ParentRepository
	ParentDetailRepo repositories.ParentDetailRepository
	ChildRepo        repositories.ChildRepository
	ObservationRepo  repositories.ObservationRepository
	Validator        Validator
	Mapper           Mapper
}

func NewDependencies(
	txRepo repositories.TransactionRepository,
	parentRepo repositories.ParentRepository,
	parentDetailRepo repositories.ParentDetailRepository,
	childRepo repositories.ChildRepository,
	observationRepo repositories.ObservationRepository,
) *Dependencies {
	return &Dependencies{
		TxRepo:           txRepo,
		ParentRepo:       parentRepo,
		ParentDetailRepo: parentDetailRepo,
		ChildRepo:        childRepo,
		ObservationRepo:  observationRepo,
		Validator:        NewRegistrationValidator(),
		Mapper:           NewRegistrationMapper(),
	}
}
