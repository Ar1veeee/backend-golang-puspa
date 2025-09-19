package observation

import (
	"backend-golang/internal/domain/repositories"
)

type Dependencies struct {
	TxRepo                   repositories.TransactionRepository
	ObservationRepo          repositories.ObservationRepository
	ObservationQuestionsRepo repositories.ObservationQuestionRepository
	ObservationAnswerRepo    repositories.ObservationAnswerRepository
	TherapistRepo            repositories.TherapistRepository
	Validator                Validator
	Mapper                   Mapper
}

func NewDependencies(
	txRepo repositories.TransactionRepository,
	observationRepo repositories.ObservationRepository,
	observationQuestionsRepo repositories.ObservationQuestionRepository,
	observationAnswerRepo repositories.ObservationAnswerRepository,
	therapistRepo repositories.TherapistRepository,
) *Dependencies {
	return &Dependencies{
		TxRepo:                   txRepo,
		ObservationRepo:          observationRepo,
		ObservationQuestionsRepo: observationQuestionsRepo,
		ObservationAnswerRepo:    observationAnswerRepo,
		TherapistRepo:            therapistRepo,
		Validator:                NewObservationValidator(),
		Mapper:                   NewObservationMapper(observationQuestionsRepo, therapistRepo),
	}
}
