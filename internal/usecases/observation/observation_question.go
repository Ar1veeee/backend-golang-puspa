package observation

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/errors"
	"context"
	"fmt"
)

type observationQuestionsUseCase struct {
	deps *Dependencies
}

func NewObservationQuestionsUseCase(deps *Dependencies) QuestionsUseCase {
	return &observationQuestionsUseCase{deps: deps}
}

func (uc *observationQuestionsUseCase) Execute(ctx context.Context, observationId int) ([]*dto.ObservationQuestionsResponse, error) {
	observation, err := uc.deps.ObservationRepo.GetById(ctx, observationId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}
	if observation == nil {
		return nil, fmt.Errorf("observation with id %d not found", observationId)
	}

	questions, err := uc.deps.ObservationQuestionsRepo.GetByAgeCategory(ctx, observation.AgeCategory)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRetrievalFailed, err)
	}

	responses := make([]*dto.ObservationQuestionsResponse, 0, len(questions))
	for _, question := range questions {
		if question == nil {
			continue
		}

		response := &dto.ObservationQuestionsResponse{
			QuestionsId:    question.Id,
			QuestionCode:   question.QuestionCode,
			AgeCategory:    question.AgeCategory,
			QuestionNumber: question.QuestionNumber,
			QuestionText:   question.QuestionText,
			Score:          question.Score,
		}
		responses = append(responses, response)
	}

	return responses, nil
}
