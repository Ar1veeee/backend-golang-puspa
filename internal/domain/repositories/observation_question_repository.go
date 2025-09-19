package repositories

import (
	"backend-golang/internal/domain/entities"
	"context"
)

type ObservationQuestionRepository interface {
	GetById(ctx context.Context, questionId int) (*entities.ObservationQuestion, error)
	GetByAgeCategory(ctx context.Context, ageCategory string) ([]*entities.ObservationQuestion, error)
}
