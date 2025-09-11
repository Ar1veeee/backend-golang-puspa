package usecase

import (
	"backend-golang/internal/registration/delivery/http/dto"
	"backend-golang/internal/registration/entity"
	"backend-golang/shared/config"
	"backend-golang/shared/constants"
	"backend-golang/shared/helpers"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type RegistrationMapper interface {
	createRequestToRegistration(req *dto.RegistrationRequest) (*entity.Parent, *entity.ParentDetail, *entity.Children, *entity.Observation, error)
}
type registrationMapper struct{}

func NewRegistrationMapper() RegistrationMapper {
	return &registrationMapper{}
}

func getEncryptionKey() string {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}
	return key
}

func stringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (m *registrationMapper) createRequestToRegistration(req *dto.RegistrationRequest) (*entity.Parent, *entity.ParentDetail, *entity.Children, *entity.Observation, error) {
	var encryptionKey = getEncryptionKey()

	parentID := helpers.GenerateULID()
	parentDetailID := helpers.GenerateULID()
	childID := helpers.GenerateULID()

	parent := &entity.Parent{
		Id:                 parentID,
		TempEmail:          req.Email,
		RegistrationStatus: string(constants.RegistrationStatusPending),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	phoneEncrypted, err := helpers.EncryptData([]byte(req.ParentPhone), encryptionKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to encrypt contact: %w", err)
	}

	addressEncrypted, err := helpers.EncryptData([]byte(req.ChildAddress), encryptionKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to encrypt address: %w", err)
	}

	parentDetail := &entity.ParentDetail{
		Id:          parentDetailID,
		ParentId:    parentID,
		ParentType:  req.ParentType,
		ParentName:  req.ParentName,
		ParentPhone: phoneEncrypted,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	child := &entity.Children{
		Id:                 childID,
		ParentId:           parentID,
		ChildName:          req.ChildName,
		ChildGender:        req.ChildGender,
		ChildBirthPlace:    req.ChildBirthPlace,
		ChildBirthDate:     req.ChildBirthDate,
		ChildAge:           req.ChildAge,
		ChildAddress:       addressEncrypted,
		ChildComplaint:     req.ChildComplaint,
		ChildSchool:        stringToPointer(req.ChildSchool),
		ChildServiceChoice: req.ChildServiceChoice,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	var ageCategory string
	switch {
	case req.ChildAge >= 2 && req.ChildAge <= 5:
		ageCategory = "Balita"
	case req.ChildAge >= 6 && req.ChildAge <= 12:
		ageCategory = "Anak-anak"
	case req.ChildAge >= 13 && req.ChildAge <= 17:
		ageCategory = "Remaja"
	default:
		ageCategory = "Lainnya"
	}

	currentTime := time.Now()
	scheduledDate := currentTime.Add(48 * time.Hour).Truncate(24 * time.Hour)
	observation := &entity.Observation{
		ChildId:       childID,
		Status:        string(constants.ObservationStatusPending),
		AgeCategory:   ageCategory,
		ScheduledDate: scheduledDate.Format("2006-01-02"),
	}

	return parent, parentDetail, child, observation, nil
}
