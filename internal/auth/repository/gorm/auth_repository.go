package gorm

import (
	"backend-golang/internal/auth/entity"
	"backend-golang/internal/auth/repository"
	"backend-golang/shared/helpers"
	"backend-golang/shared/models"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repository.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if user == nil {
		return errors.New("user data cannot be nil")
	}

	dbUser := &models.User{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
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

	return r.modelToUserEntity(&dbUser), nil
}

func (r *authRepository) FindUserByIdentifier(ctx context.Context, identifier string) (*entity.User, error) {
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

	return r.modelToUserEntity(&dbUser), nil
}

func (r *authRepository) FindUserById(ctx context.Context, id string) (*entity.User, error) {
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

	return r.modelToUserEntity(&dbUser), nil
}

func (r *authRepository) FindTokenByEmail(ctx context.Context, email string) (*entity.VerificationToken, error) {
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

func (r *authRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("email cannot be empty")
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

func (r *authRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, errors.New("username cannot be empty")
	}

	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}

func (r *authRepository) UpdateUserActiveStatus(ctx context.Context, userId string, isActive bool) error {
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

func (r *authRepository) ResetUserPassword(ctx context.Context, userId, password string) error {
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

func (r *authRepository) SaveVerificationToken(ctx context.Context, code *entity.VerificationToken) error {
	if code == nil {
		return errors.New("verification token cannot be nil")
	}

	dbCode := &models.VerificationCode{
		Id:        code.Id,
		UserId:    code.UserId,
		Code:      code.Token,
		Status:    code.Status,
		ExpiresAt: code.ExpiresAt,
		CreatedAt: code.CreatedAt,
		UpdatedAt: code.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(dbCode).Error; err != nil {
		return err
	}

	return nil
}

func (r *authRepository) VerifyAccountByToken(ctx context.Context, code string) (*entity.VerificationToken, error) {
	if code == "" {
		return nil, errors.New("code cannot be empty")
	}

	var dbCode models.VerificationCode
	if err := r.db.WithContext(ctx).
		Where("code = ? AND status = ?", code, "pending").
		First(&dbCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("verification token not found")
		}
		return nil, errors.New("failed to find verification token")
	}

	claims, err := helpers.VerifyVerificationToken(code)
	if err != nil {
		r.db.WithContext(ctx).Model(&dbCode).Update("status", "revoked")
		return nil, errors.New("invalid or expired verification token")
	}

	if claims.Subject != dbCode.UserId {
		r.db.WithContext(ctx).Model(&dbCode).Update("status", "revoked")
		return nil, errors.New("token user mismatch")
	}

	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&dbCode).
		Update("status", "used").Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update verification code: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit verification code: %w", err)
	}

	return r.modelToVerificationCodeEntity(&dbCode), nil
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
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

func (r *authRepository) FindRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
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

func (r *authRepository) RevokeRefreshToken(ctx context.Context, token string) error {
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

func (r *authRepository) GetParentByTempEmail(ctx context.Context, email string) (*entity.Parent, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var dbParent models.Parent
	if err := r.db.WithContext(ctx).Where("temp_email = ?", email).First(&dbParent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("failed to find parent by temp email")
	}

	return r.modelToParentEntity(&dbParent), nil
}

func (r *authRepository) UpdateParentUserId(ctx context.Context, tempEmail string, userId string) error {
	if tempEmail == "" || userId == "" {
		return errors.New("user_id or temp_email cannot be empty")
	}

	result := r.db.WithContext(ctx).Model(&models.Parent{}).
		Where("temp_email = ?", tempEmail).
		Update("user_id", userId)

	if result.Error != nil {
		return fmt.Errorf("failed to update parent user_id: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("parent with temp_email %s not found", tempEmail)
	}

	return nil
}

func (r *authRepository) modelToUserEntity(dbUser *models.User) *entity.User {
	return &entity.User{
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

func (r *authRepository) modelToVerificationCodeEntity(dbCode *models.VerificationCode) *entity.VerificationToken {
	return &entity.VerificationToken{
		Id:        dbCode.Id,
		UserId:    dbCode.UserId,
		Token:     dbCode.Code,
		Status:    dbCode.Status,
		ExpiresAt: dbCode.ExpiresAt,
		CreatedAt: dbCode.CreatedAt,
		UpdatedAt: dbCode.UpdatedAt,
	}
}

func (r *authRepository) modelToRefreshTokenEntity(dbToken *models.RefreshToken) *entity.RefreshToken {
	return &entity.RefreshToken{
		Id:        dbToken.Id,
		UserId:    dbToken.UserId,
		Token:     dbToken.Token,
		ExpiresAt: dbToken.ExpiresAt,
		Revoked:   dbToken.Revoked,
		CreatedAt: dbToken.CreatedAt,
	}
}

func (r *authRepository) modelToParentEntity(dbParent *models.Parent) *entity.Parent {
	var userID *string
	if dbParent.UserId != nil {
		userID = dbParent.UserId
	}

	return &entity.Parent{
		Id:        dbParent.Id,
		UserId:    userID,
		TempEmail: dbParent.TempEmail,
		CreatedAt: dbParent.CreatedAt,
		UpdatedAt: dbParent.UpdatedAt,
	}
}
