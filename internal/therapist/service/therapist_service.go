package service

import (
	"backend-golang/internal/therapist/delivery/http/dto"
	therapistErrors "backend-golang/internal/therapist/errors"
	"backend-golang/internal/therapist/repository"
	"backend-golang/shared/config"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/models"
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type TherapistService interface {
	CreateTherapist(ctx context.Context, req *dto.TherapistCreateRequest) error
	GetAllTherapist(ctx context.Context) ([]*dto.TherapistResponse, error)
}

type therapistService struct {
	therapistRepo repository.TherapistRepository
	validator     *therapistValidator
	mapper        *therapistMapper
}

func NewTherapistService(therapistRepo repository.TherapistRepository) TherapistService {
	return &therapistService{
		therapistRepo: therapistRepo,
		validator:     newTherapistValidator(),
		mapper:        newTherapistMapper(),
	}
}

func (s *therapistService) CreateTherapist(ctx context.Context, req *dto.TherapistCreateRequest) error {
	if err := s.validator.validateCreateRequest(req); err != nil {
		return err
	}

	exists, err := s.therapistRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return globalErrors.ErrEmailExists
	}

	exists, err = s.therapistRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return globalErrors.ErrUsernameExists
	}

	tx := s.therapistRepo.BeginTransaction(ctx)
	if tx == nil {
		return fmt.Errorf("%w: failed to begin transaction", globalErrors.ErrDatabaseConnection)
	}

	user, therapist, err := s.mapper.createRequestToUserAndTherapist(req)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := s.therapistRepo.CreateUserWithTx(ctx, tx, user); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", therapistErrors.ErrTherapistCreationFailed, err)
	}

	if err := s.therapistRepo.CreateTherapistWithTx(ctx, tx, therapist); err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", therapistErrors.ErrTherapistCreationFailed, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}

	return nil
}

func (s *therapistService) GetAllTherapist(ctx context.Context) ([]*dto.TherapistResponse, error) {
	therapists, err := s.therapistRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", therapistErrors.ErrTherapistRetrievalFailed, err)
	}

	if therapists == nil {
		return []*dto.TherapistResponse{}, nil
	}

	userResponses := make([]*dto.TherapistResponse, 0, len(therapists))
	for _, therapist := range therapists {
		user := therapist.User
		response := s.mapper.userToResponse(&user, therapist)
		userResponses = append(userResponses, response)
	}

	return userResponses, nil
}

func getEncryptionKey() string {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}
	return key
}

type therapistValidator struct{}

func newTherapistValidator() *therapistValidator {
	return &therapistValidator{}
}

func (v *therapistValidator) validateCreateRequest(req *dto.TherapistCreateRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	return helpers.IsValidPassword(req.Password)
}

func (v *therapistValidator) validateUpdateRequest(req *dto.TherapistUpdateRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	if req.Password != "" {
		return helpers.IsValidPassword(req.Password)
	}

	if req.Role != "" {

	}

	return nil
}

type therapistMapper struct{}

func newTherapistMapper() *therapistMapper {
	return &therapistMapper{}
}

func (m *therapistMapper) createRequestToUserAndTherapist(req *dto.TherapistCreateRequest) (*models.User, *models.Therapist, error) {
	var encryptionKey = getEncryptionKey()

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	userID := helpers.GenerateULID()
	therapistID := helpers.GenerateULID()

	user := &models.User{
		Id:        userID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	phoneEncrypted, err := helpers.EncryptData([]byte(req.TherapistPhone), encryptionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt contact: %w", err)
	}

	therapist := &models.Therapist{
		Id:               therapistID,
		UserId:           userID,
		TherapistName:    req.TherapistName,
		TherapistSection: req.TherapistSection,
		TherapistPhone:   phoneEncrypted,
	}

	return user, therapist, nil
}

func (m *therapistMapper) userToResponse(user *models.User, therapist *models.Therapist) *dto.TherapistResponse {
	var decryptedPhone string
	var decryptedKey = getEncryptionKey()

	if len(therapist.TherapistPhone) > 0 {
		decrypted, err := helpers.DecryptData(therapist.TherapistPhone, decryptedKey)
		if err == nil {
			decryptedPhone = string(decrypted)
		} else {
			decryptedPhone = ""
		}
	}

	return &dto.TherapistResponse{
		Id:               user.Id,
		Username:         user.Username,
		Email:            user.Email,
		Role:             user.Role,
		TherapistId:      therapist.Id,
		TherapistName:    therapist.TherapistName,
		TherapistSection: therapist.TherapistSection,
		TherapistPhone:   decryptedPhone,
		CreatedAt:        therapist.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        therapist.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// func (m *therapistMapper) updateRequestToUser(user *models.User, req *dto.TherapistUpdateRequest) error {
// 	if req.Username != "" {
// 		user.Username = req.Username
// 	}
// 	if req.Email != "" {
// 		user.Email = req.Email
// 	}
// 	if req.Role != "" {
// 		user.Role = req.Role
// 	}
// 	if req.Password != "" {
// 		hashedPassword, err := helpers.HashPassword(req.Password)
// 		if err != nil {
// 			return err
// 		}
// 		user.Password = hashedPassword
// 	}
// 	user.UpdatedAt = time.Now()

// 	return nil
// }
