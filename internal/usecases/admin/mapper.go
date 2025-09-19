package admin

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
	CreateRequestToUserAndAdmin(req *dto.AdminCreateRequest) (*entities.User, *entities.Admin, error)
	AdminsResponse(user *entities.User, admin *entities.Admin) (*dto.AdminResponse, error)
	UpdateRequestToUserAndAdmin(req *dto.AdminUpdateRequest, existing *entities.Admin) (*entities.User, *entities.Admin, error)
}

type adminMapper struct {
	encryptionKey string
}

func NewAdminMapper() Mapper {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}

	return &adminMapper{
		encryptionKey: key,
	}
}

func (m *adminMapper) CreateRequestToUserAndAdmin(req *dto.AdminCreateRequest) (*entities.User, *entities.Admin, error) {
	hashedPassword, err := helpers2.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	userId := helpers2.GenerateULID()
	adminId := helpers2.GenerateULID()

	user := &entities.User{
		Id:        userId,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      string(constants.RoleAdmin),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	phoneEncrypted, err := helpers2.EncryptData([]byte(req.AdminPhone), m.encryptionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt contact: %w", err)
	}

	admin := &entities.Admin{
		Id:         adminId,
		UserId:     userId,
		AdminName:  req.AdminName,
		AdminPhone: phoneEncrypted,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return user, admin, nil
}

func (m *adminMapper) AdminsResponse(user *entities.User, admin *entities.Admin) (*dto.AdminResponse, error) {
	var decryptedPhone string

	if len(admin.AdminPhone) > 0 {
		decrypted, err := helpers2.DecryptData(admin.AdminPhone, m.encryptionKey)
		if err == nil {
			decryptedPhone = string(decrypted)
		} else {
			decryptedPhone = ""
		}
	}

	return &dto.AdminResponse{
		AdminId:    admin.Id,
		AdminName:  admin.AdminName,
		Username:   user.Username,
		Email:      user.Email,
		AdminPhone: decryptedPhone,
		CreatedAt:  admin.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  admin.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (m *adminMapper) UpdateRequestToUserAndAdmin(req *dto.AdminUpdateRequest, existing *entities.Admin) (*entities.User, *entities.Admin, error) {
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

	updatedAdmin := &entities.Admin{
		Id:         existing.Id,
		UserId:     existing.UserId,
		AdminName:  existing.AdminName,
		AdminPhone: existing.AdminPhone,
		CreatedAt:  existing.CreatedAt,
		UpdatedAt:  time.Now(),
	}

	if req.AdminName != "" {
		updatedAdmin.AdminName = req.AdminName
	}
	if req.AdminPhone != "" {
		phoneEncrypted, err := helpers2.EncryptData([]byte(req.AdminPhone), m.encryptionKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt phone: %w", err)
		}
		updatedAdmin.AdminPhone = phoneEncrypted
	}

	return updatedUser, updatedAdmin, nil
}
