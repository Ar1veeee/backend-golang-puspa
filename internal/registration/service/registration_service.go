package service

import (
	"backend-golang/internal/registration/delivery/http/dto"
	registrationErrors "backend-golang/internal/registration/errors"
	"backend-golang/internal/registration/repository"
	"backend-golang/shared/config"
	"backend-golang/shared/constants"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/models"
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type RegistrationService interface {
	Registration(ctx context.Context, req *dto.RegistrationRequest) error
}

type registrationService struct {
	registrationRepo repository.RegistrationRepository
	validator        *registrationValidator
	mapper           *registrationMapper
}

func NewRegistrationService(
	registrationRepo repository.RegistrationRepository,
) RegistrationService {
	return &registrationService{
		registrationRepo: registrationRepo,
		validator:        newRegistrationValidator(),
		mapper:           newRegistrationMapper(),
	}
}

func (s *registrationService) Registration(ctx context.Context, req *dto.RegistrationRequest) error {
	if err := s.validator.validateRegisterRequest(req); err != nil {
		return err
	}

	tx := s.registrationRepo.BeginTransaction(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", globalErrors.ErrDatabaseConnection)
	}

	exists, err := s.registrationRepo.ExistsByEmail(ctx, tx, req.Email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		tx.Rollback()
		return globalErrors.ErrEmailExists
	}

	parent, parentDetail, child, observation, err := s.mapper.createRequestToRegistration(req)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := s.registrationRepo.CreateParentWithTx(ctx, tx, parent); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := s.registrationRepo.CreateParentDetailWithTx(ctx, tx, parentDetail); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := s.registrationRepo.CreateChildWithTx(ctx, tx, child); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := s.registrationRepo.CreateObservationWithTx(ctx, tx, observation); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", registrationErrors.ErrRegistrationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}

	return nil
}

type registrationValidator struct{}

func newRegistrationValidator() *registrationValidator {
	return &registrationValidator{}
}

func (v *registrationValidator) validateRegisterRequest(req *dto.RegistrationRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	return nil
}

type registrationMapper struct{}

func newRegistrationMapper() *registrationMapper {
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

func (m *registrationMapper) createRequestToRegistration(req *dto.RegistrationRequest) (*models.Parent, *models.ParentDetail, *models.Children, *models.Observation, error) {
	var encryptionKey = getEncryptionKey()

	parentID := helpers.GenerateULID()
	parentDetailID := helpers.GenerateULID()
	childID := helpers.GenerateULID()

	parent := &models.Parent{
		Id:                 parentID,
		TempEmail:          req.Email,
		RegistrationStatus: string(constants.RegistrationStatusPending),
	}

	phoneEncrypted, err := helpers.EncryptData([]byte(req.ParentPhone), encryptionKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to encrypt contact: %w", err)
	}

	addressEncrypted, err := helpers.EncryptData([]byte(req.ChildAddress), encryptionKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to encrypt address: %w", err)
	}

	parentDetail := &models.ParentDetail{
		Id:          parentDetailID,
		ParentId:    parentID,
		ParentType:  req.ParentType,
		ParentName:  req.ParentName,
		ParentPhone: phoneEncrypted,
	}

	if req.ChildAge < 0 {
		return nil, nil, nil, nil, fmt.Errorf("child age cannot be negative")
	}

	child := &models.Children{
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
	observation := &models.Observation{
		ChildId:       childID,
		Status:        string(constants.ObservationStatusPending),
		AgeCategory:   ageCategory,
		ScheduledDate: scheduledDate.Format("2006-01-02"),
	}

	return parent, parentDetail, child, observation, nil
}
