package usecase

import (
	"backend-golang/internal/observation/delivery/http/dto"
	"backend-golang/internal/observation/entity"
	"backend-golang/shared/config"
	"backend-golang/shared/helpers"
	"fmt"

	"github.com/rs/zerolog/log"
)

type ObservationMapper interface {
	ObservationsResponse(parentDetail *entity.ParentDetail, child *entity.Children, observation *entity.Observation) (*dto.ObservationsResponse, error)
	ObservationDetailResponse(
		parent *entity.Parent,
		parentDetail *entity.ParentDetail,
		child *entity.Children,
		observation *entity.Observation,
	) (*dto.DetailObservationResponse, error)
}

type observationMapper struct{}

func NewObservationMapper() ObservationMapper {
	return &observationMapper{}
}

func getEncryptionKey() string {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}
	return key
}

func (m *observationMapper) ObservationsResponse(parentDetail *entity.ParentDetail, child *entity.Children, observation *entity.Observation) (*dto.ObservationsResponse, error) {
	if observation == nil {
		return nil, fmt.Errorf("observation cannot be nil")
	}

	var parentPhone, parentName string
	var childName, childComplaint string
	var childSchool *string
	var childAge int

	var decryptionKey = getEncryptionKey()

	if parentDetail != nil {
		parentName = parentDetail.ParentName

		if len(parentDetail.ParentPhone) > 0 {
			if decryptedPhone, err := helpers.DecryptData(parentDetail.ParentPhone, decryptionKey); err != nil {
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
	parent *entity.Parent,
	parentDetail *entity.ParentDetail,
	child *entity.Children,
	observation *entity.Observation,
) (*dto.DetailObservationResponse, error) {
	if observation == nil {
		return nil, fmt.Errorf("observation cannot be nil")
	}

	var parentPhone, parentName, parentType, parentEmail string
	var childName, childComplaint, childAddress, childGender string
	var childSchool *string
	var childAge int
	var childBirthDate helpers.DateOnly

	var decryptionKey = getEncryptionKey()

	if parent != nil {
		parentEmail = parent.TempEmail
	}

	if parentDetail != nil {
		parentName = parentDetail.ParentName
		parentType = parentDetail.ParentType

		if len(parentDetail.ParentPhone) > 0 {
			if decryptedPhone, err := helpers.DecryptData(parentDetail.ParentPhone, decryptionKey); err != nil {
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
			if decryptedAddress, err := helpers.DecryptData(child.ChildAddress, decryptionKey); err != nil {
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
