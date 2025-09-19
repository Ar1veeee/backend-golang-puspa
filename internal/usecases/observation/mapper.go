package observation

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/helpers"
	"backend-golang/internal/infrastructure/config"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type Mapper interface {
	ObservationsResponse(
		parentDetail *entities.ParentDetail,
		child *entities.Children,
		observation *entities.Observation,
	) (*dto.ObservationsResponse, error)
	ObservationDetailResponse(
		parent *entities.Parent,
		parentDetail *entities.ParentDetail,
		child *entities.Children,
		observation *entities.Observation,
	) (*dto.DetailObservationResponse, error)
	UpdateToObservationAndCreateToAnswer(ctx context.Context, observationId int, req *dto.SubmitObservationRequest) (*entities.Observation, []*entities.ObservationAnswer, error)
}

type observationMapper struct {
	encryptionKey            string
	observationQuestionsRepo repositories.ObservationQuestionRepository
	therapistRepo            repositories.TherapistRepository
}

func NewObservationMapper(
	observationQuestionsRepo repositories.ObservationQuestionRepository,
	therapistRepo repositories.TherapistRepository,
) Mapper {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}

	return &observationMapper{
		encryptionKey:            key,
		observationQuestionsRepo: observationQuestionsRepo,
		therapistRepo:            therapistRepo,
	}
}

func (m *observationMapper) ObservationsResponse(parentDetail *entities.ParentDetail, child *entities.Children, observation *entities.Observation) (*dto.ObservationsResponse, error) {
	if observation == nil {
		return nil, fmt.Errorf("observation cannot be nil")
	}

	var parentPhone, parentName string
	var childName, childComplaint string
	var childSchool *string
	var childAge int

	if parentDetail != nil {
		parentName = parentDetail.ParentName

		if len(parentDetail.ParentPhone) > 0 {
			if decryptedPhone, err := helpers.DecryptData(parentDetail.ParentPhone, m.encryptionKey); err != nil {
				log.Warn().Err(err).Msg("Failed to decrypt parent phone")
				parentPhone = "[Encrypted]"
			} else {
				parentPhone = string(decryptedPhone)
			}
		}
	}

	if child != nil {
		childName = child.ChildName
		childSchool = child.ChildSchool
		childComplaint = child.ChildComplaint

		if !child.ChildBirthDate.ToTime().IsZero() {
			birthTime := child.ChildBirthDate.ToTime()
			childAge = helpers.CalculateAge(birthTime)
		}
	}

	return &dto.ObservationsResponse{
		ObservationId:  observation.Id,
		AgeCategory:    observation.AgeCategory,
		ChildName:      childName,
		ChildSchool:    childSchool,
		ChildAge:       childAge,
		ChildComplaint: childComplaint,
		ParentName:     parentName,
		ParentPhone:    parentPhone,
		ScheduledDate:  observation.ScheduledDate,
		Status:         observation.Status,
	}, nil
}

func (m *observationMapper) ObservationDetailResponse(
	parent *entities.Parent,
	parentDetail *entities.ParentDetail,
	child *entities.Children,
	observation *entities.Observation,
) (*dto.DetailObservationResponse, error) {
	if observation == nil {
		return nil, fmt.Errorf("observation cannot be nil")
	}

	var parentPhone, parentName, parentType, parentEmail string
	var childName, childComplaint, childAddress, childGender string
	var childSchool *string
	var childAge int
	var childBirthDate helpers.DateOnly

	if parent != nil {
		parentEmail = parent.TempEmail
	}

	if parentDetail != nil {
		parentName = parentDetail.ParentName
		parentType = parentDetail.ParentType

		if len(parentDetail.ParentPhone) > 0 {
			if decryptedPhone, err := helpers.DecryptData(parentDetail.ParentPhone, m.encryptionKey); err != nil {
				log.Warn().Err(err).Msg("Failed to decrypt parent phone")
				parentPhone = "[Encrypted]"
			} else {
				parentPhone = string(decryptedPhone)
			}
		}
	}

	if child != nil {
		childName = child.ChildName
		childSchool = child.ChildSchool
		childComplaint = child.ChildComplaint
		childGender = child.ChildGender

		if len(child.ChildAddress) > 0 {
			if decryptedAddress, err := helpers.DecryptData(child.ChildAddress, m.encryptionKey); err != nil {
				log.Warn().Err(err).Msg("Failed to decrypt parent phone")
				childAddress = "[Encrypted]"
			} else {
				childAddress = string(decryptedAddress)
			}
		}

		if !child.ChildBirthDate.ToTime().IsZero() {
			birthTime := child.ChildBirthDate.ToTime()
			childBirthDate = child.ChildBirthDate
			childAge = helpers.CalculateAge(birthTime)
		}
	}

	return &dto.DetailObservationResponse{
		ObservationId:  observation.Id,
		ChildName:      childName,
		ChildBirthDate: childBirthDate,
		ChildAge:       childAge,
		ChildGender:    childGender,
		ChildSchool:    childSchool,
		ChildAddress:   childAddress,
		ParentName:     parentName,
		ParentType:     parentType,
		ParentPhone:    parentPhone,
		ChildComplaint: childComplaint,
		Email:          parentEmail,
	}, nil
}

func (m *observationMapper) UpdateToObservationAndCreateToAnswer(ctx context.Context, observationId int, req *dto.SubmitObservationRequest) (*entities.Observation, []*entities.ObservationAnswer, error) {
	userId, ok := helpers.GetUserID(ctx)
	if !ok {
		return nil, nil, fmt.Errorf("userId not found in context")
	}

	therapist, err := m.therapistRepo.GetByUserId(ctx, userId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get therapist: %w", err)
	}

	if req == nil {
		return nil, nil, errors.New("request cannot be nil")
	}

	observation := &entities.Observation{
		TherapistId:    therapist.Id,
		Conclusion:     req.Conclusion,
		Recommendation: req.Recommendation,
		UpdatedAt:      time.Now(),
	}

	if len(req.Answers) == 0 {
		return observation, nil, errors.New("at least one answer is required")
	}

	observationAnswers := make([]*entities.ObservationAnswer, 0, len(req.Answers))
	var totalScore int

	for _, answerInput := range req.Answers {
		question, err := m.observationQuestionsRepo.GetById(ctx, answerInput.QuestionId)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get question with id %d: %w", answerInput.QuestionId, err)
		}
		if question == nil {
			return nil, nil, fmt.Errorf("question with id %d not found", answerInput.QuestionId)
		}

		var scoreEarned int
		if answerInput.Answer {
			scoreEarned = question.Score
		} else {
			scoreEarned = 0
		}

		var note *string
		if answerInput.Note != "" {
			note = &answerInput.Note
		}

		observationAnswer := &entities.ObservationAnswer{
			ObservationId: observationId,
			QuestionId:    answerInput.QuestionId,
			Answer:        answerInput.Answer,
			ScoreEarned:   scoreEarned,
			Note:          note,
		}

		observationAnswers = append(observationAnswers, observationAnswer)
		totalScore += scoreEarned
		log.Debug().Int("total_score", totalScore).Msg("Calculated total score")
	}

	observation.TotalScore = totalScore

	return observation, observationAnswers, nil
}
