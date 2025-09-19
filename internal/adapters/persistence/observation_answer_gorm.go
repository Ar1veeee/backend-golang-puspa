package persistence

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/infrastructure/database/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type observationAnswerRepository struct {
	db *gorm.DB
}

func NewObservationAnswerRepository(db *gorm.DB) repositories.ObservationAnswerRepository {
	return &observationAnswerRepository{
		db: db,
	}
}

func (r *observationAnswerRepository) Create(ctx context.Context, tx *gorm.DB, answers []*entities.ObservationAnswer) error {
	if len(answers) == 0 {
		return errors.New("answer cannot be empty")
	}

	dbAnswers := make([]*models.ObservationAnswer, 0, len(answers))
	for _, answer := range answers {
		if answer == nil {
			return errors.New("answer cannot be nil")
		}
		dbAnswer := r.entityToModel(answer)
		dbAnswers = append(dbAnswers, dbAnswer)
	}

	if err := tx.WithContext(ctx).Create(dbAnswers).Error; err != nil {
		return fmt.Errorf("failed to create observation_answers: %w", err)
	}

	return nil
}

func (r *observationAnswerRepository) entityToModel(answer *entities.ObservationAnswer) *models.ObservationAnswer {
	return &models.ObservationAnswer{
		Id:            answer.Id,
		ObservationId: answer.ObservationId,
		QuestionId:    answer.QuestionId,
		Answer:        answer.Answer,
		ScoreEarned:   answer.ScoreEarned,
		Note:          answer.Note,
	}
}
