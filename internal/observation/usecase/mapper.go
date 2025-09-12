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
	AllObservationsResponse(parentDetail *entity.ParentDetail, child *entity.Children, observation *entity.Observation) (*dto.PendingObservationsResponse, error)
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

func (m *observationMapper) AllObservationsResponse(parentDetail *entity.ParentDetail, child *entity.Children, observation *entity.Observation) (*dto.PendingObservationsResponse, error) {
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

	return &dto.PendingObservationsResponse{
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
