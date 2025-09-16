package persistence

import (
	"backend-golang/internal/constants"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/internal/helpers"
	"context"

	"backend-golang/pkg/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type verificationTokenRepository struct {
	db *gorm.DB
}

func NewVerificationTokenRepository(db *gorm.DB) repositories.VerificationTokenRepository {
	return &verificationTokenRepository{db: db}
}

func (r *verificationTokenRepository) Create(ctx context.Context, token *entities.VerificationToken) error {
	if token == nil {
		return errors.New("verification token cannot be nil")
	}

	dbCode := &models.VerificationCode{
		Id:        token.Id,
		UserId:    token.UserId,
		Code:      token.Token,
		Status:    token.Status,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(dbCode).Error; err != nil {
		return err
	}

	return nil
}

func (r *verificationTokenRepository) GetByEmail(ctx context.Context, email string) (*entities.VerificationToken, error) {
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

	var dbToken models.VerificationCode
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", dbUser.Id, "pending").
		Order("created_at DESC").
		First(&dbToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("failed to find user by email")
	}

	if dbToken.ExpiresAt.Before(time.Now()) {
		return nil, nil
	}

	return r.modelToVerificationCodeEntity(&dbToken), nil
}

func (r *verificationTokenRepository) GetByToken(ctx context.Context, code string) (*entities.VerificationToken, error) {

	if code == "" {
		return nil, errors.New("code cannot be empty")
	}

	var dbCode models.VerificationCode
	if err := r.db.WithContext(ctx).
		Where("code = ? AND status = ?", code, string(constants.VerificationCodeStatusPending)).
		First(&dbCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("verification token not found")
		}
		return nil, errors.New("failed to find verification token")
	}

	claims, err := helpers.VerifyVerificationToken(code)
	if err != nil {
		r.db.WithContext(ctx).Model(&dbCode).Update("status", string(constants.VerificationCodeStatusRevoked))
		return nil, errors.New("invalid or expired verification token")
	}

	if claims.Subject != dbCode.UserId {
		return nil, errors.New("token user mismatch")
	}

	return r.modelToVerificationCodeEntity(&dbCode), nil
}

func (r *verificationTokenRepository) UpdateStatus(ctx context.Context, code string) error {
	if code == "" {
		return errors.New("code cannot be empty")
	}

	var dbCode models.VerificationCode
	if err := r.db.WithContext(ctx).
		Where("code = ? AND status = ?", code, string(constants.VerificationCodeStatusPending)).
		First(&dbCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("verification token not found")
		}
		return errors.New("failed to find verification token")
	}

	claims, err := helpers.VerifyVerificationToken(code)
	if err != nil {
		r.db.WithContext(ctx).Model(&dbCode).Update("status", string(constants.VerificationCodeStatusRevoked))
		return errors.New("invalid or expired verification token")
	}

	if claims.Subject != dbCode.UserId {
		r.db.WithContext(ctx).Model(&dbCode).Update("status", string(constants.VerificationCodeStatusRevoked))
		return errors.New("token user mismatch")
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&dbCode).
		Update("status", string(constants.VerificationCodeStatusUsed)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update verification code: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit verification code: %w", err)
	}

	return nil
}

func (r *verificationTokenRepository) modelToVerificationCodeEntity(dbCode *models.VerificationCode) *entities.VerificationToken {
	return &entities.VerificationToken{
		Id:        dbCode.Id,
		UserId:    dbCode.UserId,
		Token:     dbCode.Code,
		Status:    dbCode.Status,
		ExpiresAt: dbCode.ExpiresAt,
		CreatedAt: dbCode.CreatedAt,
		UpdatedAt: dbCode.UpdatedAt,
	}
}
