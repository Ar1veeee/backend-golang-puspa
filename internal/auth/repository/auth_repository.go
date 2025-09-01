package repository

import (
	"backend-golang/shared/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

type AuthRepository interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, username string, user *models.User) error
	RefreshToken(ctx context.Context, refreshToken string, token *models.RefreshToken) error
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error
	Logout(ctx context.Context, refreshToken string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *authRepository) Login(ctx context.Context, username string, user *models.User) error {
	if username == "" {
		return errors.New("username cannot be nil")
	}

	result := r.db.WithContext(ctx).Where("username = ?", username).First(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *authRepository) RefreshToken(ctx context.Context, refreshToken string, token *models.RefreshToken) error {
	if refreshToken == "" {
		return errors.New("refresh token cannot be nil")
	}

	result := r.db.WithContext(ctx).Where("token = ? AND revoked = ?", refreshToken, false).First(&token)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	if token == nil {
		return errors.New("token cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(token)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *authRepository) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return errors.New("refresh token cannot be nil")
	}

	result := r.db.WithContext(ctx).Model(&models.RefreshToken{}).Where("token = ? AND revoked = ?", refreshToken, false).Update("revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *authRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

func (r *authRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}
