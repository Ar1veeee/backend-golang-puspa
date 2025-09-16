package mapper

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

type TherapistMapper interface {
	CreateRequestToUserAndTherapist(req *dto.TherapistCreateRequest) (*entities.User, *entities.Therapist, error)
	TherapistsResponse(user *entities.User, therapist *entities.Therapist) (*dto.TherapistResponse, error)
	UpdateRequestToUserAndTherapist(req *dto.TherapistUpdateRequest, existing *entities.Therapist) (*entities.User, *entities.Therapist, error)
}

type therapistMapper struct {
	encryptionKey string
}

func NewTherapistMapper() TherapistMapper {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}

	return &therapistMapper{
		encryptionKey: key,
	}
}

func (m *therapistMapper) CreateRequestToUserAndTherapist(req *dto.TherapistCreateRequest) (*entities.User, *entities.Therapist, error) {
	hashedPassword, err := helpers2.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	userId := helpers2.GenerateULID()
	therapistID := helpers2.GenerateULID()

	user := &entities.User{
		Id:        userId,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      string(constants.RoleTherapist),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	phoneEncrypted, err := helpers2.EncryptData([]byte(req.TherapistPhone), m.encryptionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt contact: %w", err)
	}

	therapist := &entities.Therapist{
		Id:               therapistID,
		UserId:           userId,
		TherapistName:    req.TherapistName,
		TherapistSection: req.TherapistSection,
		TherapistPhone:   phoneEncrypted,
	}

	return user, therapist, nil
}

func (m *therapistMapper) TherapistsResponse(user *entities.User, therapist *entities.Therapist) (*dto.TherapistResponse, error) {
	var decryptedPhone string

	if len(therapist.TherapistPhone) > 0 {
		decrypted, err := helpers2.DecryptData(therapist.TherapistPhone, m.encryptionKey)
		if err == nil {
			decryptedPhone = string(decrypted)
		} else {
			decryptedPhone = ""
		}
	}

	return &dto.TherapistResponse{
		UserId:           user.Id,
		TherapistId:      therapist.Id,
		TherapistName:    therapist.TherapistName,
		TherapistSection: therapist.TherapistSection,
		Username:         user.Username,
		Email:            user.Email,
		TherapistPhone:   decryptedPhone,
		CreatedAt:        therapist.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        therapist.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (m *therapistMapper) UpdateRequestToUserAndTherapist(req *dto.TherapistUpdateRequest, existing *entities.Therapist) (*entities.User, *entities.Therapist, error) {
	updatedUser := &entities.User{
		Id:        existing.User.Id,
		Username:  existing.User.Username,
		Email:     existing.User.Email,
		Password:  existing.User.Password,
		Role:      existing.User.Role,
		IsActive:  existing.User.IsActive,
		CreatedAt: existing.User.CreatedAt,
		UpdatedAt: time.Now(),
		LastLogin: existing.User.LastLogin,
	}

	if req.Username != "" {
		updatedUser.Username = req.Username
	}
	if req.Email != "" {
		updatedUser.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := helpers2.HashPassword(req.Password)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to hash password: %w", err)
		}
		updatedUser.Password = hashedPassword
	}

	updatedTherapist := &entities.Therapist{
		Id:               existing.Id,
		UserId:           existing.UserId,
		TherapistName:    existing.TherapistName,
		TherapistSection: existing.TherapistSection,
		TherapistPhone:   existing.TherapistPhone,
		CreatedAt:        existing.CreatedAt,
		UpdatedAt:        time.Now(),
	}

	if req.TherapistName != "" {
		updatedTherapist.TherapistName = req.TherapistName
	}

	if req.TherapistSection != "" {
		updatedTherapist.TherapistSection = req.TherapistSection
	}

	if req.TherapistPhone != "" {
		phoneEncrypted, err := helpers2.EncryptData([]byte(req.TherapistPhone), m.encryptionKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt phone: %w", err)
		}
		updatedTherapist.TherapistPhone = phoneEncrypted
	}

	return updatedUser, updatedTherapist, nil
}
