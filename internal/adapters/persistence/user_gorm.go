package persistence

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/pkg/models"
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, tx *gorm.DB, user *entities.User) error {
	dbUser := r.entityToModel(user)

	if err := tx.WithContext(ctx).Create(dbUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var dbUser models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user with this email not found")
		}
		return nil, errors.New("failed to find user by email")
	}

	return r.modelToEntity(&dbUser), nil
}

func (r *userRepository) GetByIdentifier(ctx context.Context, identifier string) (*entities.User, error) {
	if identifier == "" {
		return nil, errors.New("identifier cannot be empty")
	}

	var dbUser models.User
	if err := r.db.WithContext(ctx).
		Where("username = ? OR email = ?", identifier, identifier).
		First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user with this username or email not found")
		}
		return nil, errors.New("failed to find user by username or email")
	}

	return r.modelToEntity(&dbUser), nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*entities.User, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	var dbUser models.User
	if err := r.db.WithContext(ctx).First(&dbUser, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user with this id not found")
		}
		return nil, errors.New("failed to find user by id")
	}

	return r.modelToEntity(&dbUser), nil
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user *entities.User) error {
	if user == nil {
		return errors.New("user cannot be empty")
	}

	dbUser := r.entityToModel(user)

	if err := tx.WithContext(ctx).Save(dbUser).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) UpdateActiveStatus(ctx context.Context, userId string, isActive bool) error {
	if userId == "" {
		return errors.New("user id cannot be empty")
	}

	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"is_active":  isActive,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return errors.New("failed to update user active status")
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userId, password string) error {
	if userId == "" {
		return errors.New("user id cannot be empty")
	}

	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"password":   password,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return errors.New("failed to reset password")
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) CheckExisting(ctx context.Context, email, username string) (emailExists, usernameExists bool, err error) {
	var emailCount, usernameCount int64

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&emailCount).Error
	})

	eg.Go(func() error {
		return r.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&usernameCount).Error
	})

	if err := eg.Wait(); err != nil {
		return false, false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return emailCount > 0, usernameCount > 0, nil
}

func (r *userRepository) modelToEntity(dbUser *models.User) *entities.User {
	return &entities.User{
		Id:        dbUser.Id,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		Role:      dbUser.Role,
		IsActive:  dbUser.IsActive,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

func (r *userRepository) entityToModel(user *entities.User) *models.User {
	return &models.User{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
