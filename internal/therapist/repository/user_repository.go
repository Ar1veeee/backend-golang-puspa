package repository

import (
	"backend-golang/shared/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetAll(ctx context.Context) ([]*models.User, error)
	GetById(ctx context.Context, id string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	var users []*models.User

	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	var user models.User

	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.Id == "" {
		return errors.New("user id cannot be nil")
	}

	result := r.db.WithContext(ctx).Model(user).Select("*").Omit("createdAt").Updates(user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}
