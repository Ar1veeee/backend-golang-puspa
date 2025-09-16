package persistence

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"context"

	"backend-golang/pkg/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) repositories.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entities.RefreshToken) error {
	if token == nil {
		return errors.New("refresh token cannot be nil")
	}

	dbToken := &models.RefreshToken{
		Id:        token.Id,
		UserId:    token.UserId,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		Revoked:   token.Revoked,
		CreatedAt: token.CreatedAt,
	}

	if err := r.db.WithContext(ctx).Create(&dbToken).Error; err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	if token == "" {
		return nil, errors.New("refresh token cannot be nil")
	}

	var dbToken models.RefreshToken

	if err := r.db.WithContext(ctx).
		Where("token = ? AND revoked = ?", token, false).
		First(&dbToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, errors.New("failed to find refresh token")
	}

	return r.modelToRefreshTokenEntity(&dbToken), nil
}

func (r *refreshTokenRepository) RevokeStatus(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("refresh token cannot be empty")
	}

	result := r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("token = ? AND revoked = ?", token, false).
		Update("revoked", true)

	if result.Error != nil {
		return errors.New("failed to revoke refresh token")
	}

	if result.RowsAffected == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}

func (r *refreshTokenRepository) modelToRefreshTokenEntity(dbToken *models.RefreshToken) *entities.RefreshToken {
	return &entities.RefreshToken{
		Id:        dbToken.Id,
		UserId:    dbToken.UserId,
		Token:     dbToken.Token,
		ExpiresAt: dbToken.ExpiresAt,
		Revoked:   dbToken.Revoked,
		CreatedAt: dbToken.CreatedAt,
	}
}
