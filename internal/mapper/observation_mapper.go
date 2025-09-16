package mapper

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	helpers2 "backend-golang/internal/helpers"
	"backend-golang/internal/infrastructure/config"
	"fmt"

	"github.com/rs/zerolog/log"
)

type ObservationMapper interface {
	ObservationsResponse(parentDetail *entities.ParentDetail, child *entities.Children, observation *entities.Observation) (*dto.ObservationsResponse, error)
	ObservationDetailResponse(
		parent *entities.Parent,
		parentDetail *entities.ParentDetail,
		child *entities.Children,
		observation *entities.Observation,
	) (*dto.DetailObservationResponse, error)
}

type observationMapper struct {
	encryptionKey string
}

func NewObservationMapper() ObservationMapper {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}

	return &observationMapper{
		encryptionKey: key,
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
			if decryptedPhone, err := helpers2.DecryptData(parentDetail.ParentPhone, m.encryptionKey); err != nil {
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
			childAge = helpers2.CalculateAge(birthTime)
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
	var childBirthDate helpers2.DateOnly

	if parent != nil {
		parentEmail = parent.TempEmail
	}

	if parentDetail != nil {
		parentName = parentDetail.ParentName
		parentType = parentDetail.ParentType

		if len(parentDetail.ParentPhone) > 0 {
			if decryptedPhone, err := helpers2.DecryptData(parentDetail.ParentPhone, m.encryptionKey); err != nil {
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
			if decryptedAddress, err := helpers2.DecryptData(child.ChildAddress, m.encryptionKey); err != nil {
				log.Warn().Err(err).Msg("Failed to decrypt parent phone")
				childAddress = "[Encrypted]"
			} else {
				childAddress = string(decryptedAddress)
			}
		}

		if !child.ChildBirthDate.ToTime().IsZero() {
			birthTime := child.ChildBirthDate.ToTime()
			childBirthDate = child.ChildBirthDate
			childAge = helpers2.CalculateAge(birthTime)
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
