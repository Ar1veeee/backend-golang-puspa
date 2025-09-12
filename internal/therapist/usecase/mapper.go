package usecase

import (
	"backend-golang/internal/therapist/delivery/http/dto"
	"backend-golang/internal/therapist/entity"
	"backend-golang/shared/config"
	"backend-golang/shared/helpers"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type TherapistMapper interface {
	CreateRequestToUserAndTherapist(req *dto.TherapistCreateRequest) (*entity.User, *entity.Therapist, error)
	AllTherapistsResponse(user *entity.User, therapist *entity.Therapist) *dto.TherapistResponse
}

func getEncryptionKey() string {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}
	return key
}

type therapistMapper struct{}

func NewTherapistMapper() TherapistMapper {
	return &therapistMapper{}
}

func (m *therapistMapper) CreateRequestToUserAndTherapist(req *dto.TherapistCreateRequest) (*entity.User, *entity.Therapist, error) {
	var encryptionKey = getEncryptionKey()

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	userID := helpers.GenerateULID()
	therapistID := helpers.GenerateULID()

	user := &entity.User{
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

	therapist := &entity.Therapist{
		Id:               therapistID,
		UserId:           userID,
		TherapistName:    req.TherapistName,
		TherapistSection: req.TherapistSection,
		TherapistPhone:   phoneEncrypted,
	}

	return user, therapist, nil
}

func (m *therapistMapper) AllTherapistsResponse(user *entity.User, therapist *entity.Therapist) *dto.TherapistResponse {
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
