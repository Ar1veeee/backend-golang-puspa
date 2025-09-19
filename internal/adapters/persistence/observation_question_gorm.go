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

type observationQuestionRepository struct {
	db *gorm.DB
}

func NewObservationQuestionRepository(db *gorm.DB) repositories.ObservationQuestionRepository {
	return &observationQuestionRepository{
		db: db,
	}
}

func (r *observationQuestionRepository) GetById(ctx context.Context, questionId int) (*entities.ObservationQuestion, error) {
	var dbQuestion models.ObservationQuestion

	if err := r.db.First(&dbQuestion, "id = ?", questionId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ObservationQuestion Not Found")
		}
		return nil, err
	}

	question := r.modelToEntity(&dbQuestion)
	if question == nil {
		return nil, fmt.Errorf("ObservationQuestion Not Found")
	}

	return question, nil
}

func (r *observationQuestionRepository) GetByAgeCategory(ctx context.Context, ageCategory string) ([]*entities.ObservationQuestion, error) {
	if ageCategory == "" {
		return nil, errors.New("age category is required")
	}

	var dbObservationQuestions []*models.ObservationQuestion
	if err := r.db.WithContext(ctx).
		Where("age_category = ?", ageCategory).
		Order("question_number asc").
		Find(&dbObservationQuestions).Error; err != nil {
		return nil, fmt.Errorf("failed to get question: %w", err)
	}

	questions := make([]*entities.ObservationQuestion, 0, len(dbObservationQuestions))
	for _, dbObservationQuestion := range dbObservationQuestions {
		if dbObservationQuestion == nil {
			continue
		}
		question := r.modelToEntity(dbObservationQuestion)
		questions = append(questions, question)
	}

	return questions, nil
}

func (r *observationQuestionRepository) modelToEntity(dbObservationQuestion *models.ObservationQuestion) *entities.ObservationQuestion {
	observationQuestion := &entities.ObservationQuestion{
		Id:             dbObservationQuestion.Id,
		QuestionCode:   dbObservationQuestion.QuestionCode,
		AgeCategory:    dbObservationQuestion.AgeCategory,
		QuestionNumber: dbObservationQuestion.QuestionNumber,
		QuestionText:   dbObservationQuestion.QuestionText,
		Score:          dbObservationQuestion.Score,
		IsActive:       dbObservationQuestion.IsActive,
		CreatedAt:      dbObservationQuestion.CreatedAt,
	}

	return observationQuestion
}
