package registration

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/constants"
	"backend-golang/internal/domain/entities"
	helpers2 "backend-golang/internal/helpers"
	"backend-golang/internal/infrastructure/config"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type Mapper interface {
	CreateRequestToRegistration(req *dto.RegistrationRequest) (*entities.Parent, *entities.ParentDetail, *entities.Children, *entities.Observation, error)
}
type registrationMapper struct {
	encryptionKey string
}

func NewRegistrationMapper() Mapper {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}

	return &registrationMapper{
		encryptionKey: key,
	}
}

func (m *registrationMapper) CreateRequestToRegistration(req *dto.RegistrationRequest) (*entities.Parent, *entities.ParentDetail, *entities.Children, *entities.Observation, error) {
	parentID := helpers2.GenerateULID()
	parentDetailID := helpers2.GenerateULID()
	childID := helpers2.GenerateULID()

	parent := &entities.Parent{
		Id:                 parentID,
		TempEmail:          req.Email,
		RegistrationStatus: string(constants.RegistrationStatusPending),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	phoneEncrypted, err := helpers2.EncryptData([]byte(req.ParentPhone), m.encryptionKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to encrypt contact: %w", err)
	}

	addressEncrypted, err := helpers2.EncryptData([]byte(req.ChildAddress), m.encryptionKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to encrypt address: %w", err)
	}

	parentDetail := &entities.ParentDetail{
		Id:          parentDetailID,
		ParentId:    parentID,
		ParentType:  req.ParentType,
		ParentName:  req.ParentName,
		ParentPhone: phoneEncrypted,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	child := &entities.Children{
		Id:                 childID,
		ParentId:           parentID,
		ChildName:          req.ChildName,
		ChildGender:        req.ChildGender,
		ChildBirthPlace:    req.ChildBirthPlace,
		ChildBirthDate:     req.ChildBirthDate,
		ChildAddress:       addressEncrypted,
		ChildComplaint:     req.ChildComplaint,
		ChildSchool:        req.ChildSchool,
		ChildServiceChoice: req.ChildServiceChoice,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	var childAge int

	if !req.ChildBirthDate.ToTime().IsZero() {
		birthTime := req.ChildBirthDate.ToTime()
		childAge = helpers2.CalculateAge(birthTime)
	}

	var ageCategory string
	switch {
	case childAge >= 0 && childAge <= 5:
		ageCategory = "Balita"
	case childAge >= 6 && childAge <= 12:
		ageCategory = "Anak-anak"
	case childAge >= 13 && childAge <= 17:
		ageCategory = "Remaja"
	default:
		ageCategory = "Lainnya"
	}

	currentTime := time.Now()
	scheduledDate := currentTime.Add(48 * time.Hour).Truncate(24 * time.Hour)
	observation := &entities.Observation{
		ChildId:       childID,
		Status:        string(constants.ObservationStatusPending),
		AgeCategory:   ageCategory,
		ScheduledDate: helpers2.DateOnly(scheduledDate),
	}

	return parent, parentDetail, child, observation, nil
}
