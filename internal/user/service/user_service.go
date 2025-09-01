package service

import (
	"backend-golang/internal/user/dto"
	userErrors "backend-golang/internal/user/errors"
	"backend-golang/internal/user/repository"
	globalErrors "backend-golang/shared/errors"
	"backend-golang/shared/helpers"
	"backend-golang/shared/models"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dto.UserCreateRequest) (*dto.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]*dto.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req *dto.UserUpdateRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	userRepo  repository.UserRepository
	validator *userValidator
	mapper    *userMapper
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo:  userRepo,
		validator: newUserValidator(),
		mapper:    newUserMapper(),
	}
}

func (u *userService) CreateUser(ctx context.Context, req *dto.UserCreateRequest) (*dto.UserResponse, error) {
	if err := u.validator.validateCreateRequest(req); err != nil {
		return nil, err
	}

	exists, err := u.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return nil, userErrors.ErrEmailExists
	}

	exists, err = u.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
	}
	if exists {
		return nil, userErrors.ErrUsernameExists
	}

	user, err := u.mapper.createRequestToUser(req)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("%w: %v", userErrors.ErrUserCreationFailed, err)
	}

	return u.mapper.userToResponse(user), nil
}

func (u *userService) GetAllUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", userErrors.ErrUserRetrievalFailed, err)
	}

	userResponses := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, u.mapper.userToResponse(user))
	}

	return userResponses, nil
}

func (u *userService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	if id == "" {
		return nil, userErrors.ErrUserIDRequired
	}
	if !helpers.IsValidUUID(id) {
		return nil, userErrors.ErrInvalidUserID
	}

	user, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", userErrors.ErrUserRetrievalFailed, err)
	}
	return u.mapper.userToResponse(user), nil
}

func (u *userService) UpdateUser(ctx context.Context, id string, req *dto.UserUpdateRequest) (*dto.UserResponse, error) {
	existingUser, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", userErrors.ErrUserRetrievalFailed, err)
	}

	if req.Email != "" && req.Email != existingUser.Email {
		exists, err := u.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
		}
		if exists {
			return nil, userErrors.ErrEmailExists
		}
	}

	if req.Username != "" && req.Username != existingUser.Username {
		exists, err := u.userRepo.ExistsByUsername(ctx, req.Username)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", globalErrors.ErrDatabaseConnection, err)
		}
		if exists {
			return nil, userErrors.ErrUsernameExists
		}
	}

	if err := u.mapper.updateRequestToUser(existingUser, req); err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	if err := u.userRepo.Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("%w: %v", userErrors.ErrUserUpdateFailed, err)
	}

	return u.mapper.userToResponse(existingUser), nil
}

func (u *userService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return userErrors.ErrUserIDRequired
	}

	if !helpers.IsValidUUID(id) {
		return userErrors.ErrInvalidUserID
	}

	_, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userErrors.ErrUserNotFound
		}
		return fmt.Errorf("%w: %v", userErrors.ErrUserRetrievalFailed, err)
	}

	if err := u.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", userErrors.ErrUserDeletionFailed, err)
	}

	return nil
}

type userValidator struct{}

func newUserValidator() *userValidator {
	return &userValidator{}
}

func (v *userValidator) validateCreateRequest(req *dto.UserCreateRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	return helpers.IsValidPassword(req.Password)
}

func (v *userValidator) validateUpdateRequest(req *dto.UserUpdateRequest) error {
	if err := helpers.ValidateStruct(req); err != nil {
		return err
	}

	if req.Password != "" {
		return helpers.IsValidPassword(req.Password)
	}

	return nil
}

type userMapper struct{}

func newUserMapper() *userMapper {
	return &userMapper{}
}

func (m *userMapper) createRequestToUser(req *dto.UserCreateRequest) (*models.User, error) {
	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *userMapper) updateRequestToUser(user *models.User, req *dto.UserUpdateRequest) error {
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Password != "" {
		hashedPassword, err := helpers.HashPassword(req.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	user.UpdatedAt = time.Now()

	return nil
}

func (m *userMapper) userToResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		CreateAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
