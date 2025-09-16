package persistence

import (
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/domain/repositories"
	"backend-golang/pkg/models"
	"errors"
	"fmt"

	"context"

	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) repositories.AdminRepository {
	return &adminRepository{
		db: db,
	}
}

func (r *adminRepository) Create(ctx context.Context, tx *gorm.DB, admin *entities.Admin) error {
	if admin == nil {
		return errors.New("admin data cannot be empty")
	}

	dbAdmin := r.entityToModel(admin)

	if err := tx.WithContext(ctx).Create(dbAdmin).Error; err != nil {
		return fmt.Errorf("failed create therapist: %w", err)
	}

	return nil
}

func (r *adminRepository) GetById(ctx context.Context, adminId string) (*entities.Admin, error) {
	var dbAdmin models.Admin

	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", adminId).
		First(&dbAdmin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, fmt.Errorf("failed to find admin: %w", err)
	}

	admin := r.modelToEntity(&dbAdmin)
	if admin == nil {
		return nil, errors.New("failed to find admin by id")
	}

	return admin, nil
}

func (r *adminRepository) GetAll(ctx context.Context) ([]*entities.Admin, error) {
	var dbAdmins []*models.Admin

	if err := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = admins.user_id").
		Where("users.is_active = ?", true).
		Preload("User").
		Find(&dbAdmins).Error; err != nil {
		return nil, fmt.Errorf("failed find admins: %w", err)
	}

	admins := make([]*entities.Admin, 0, len(dbAdmins))
	for _, dbAdmin := range dbAdmins {
		admin := r.modelToEntity(dbAdmin)
		admins = append(admins, admin)
	}

	return admins, nil
}

func (r *adminRepository) Update(ctx context.Context, tx *gorm.DB, admin *entities.Admin) error {
	if admin == nil {
		return errors.New("admin data cannot be empty")
	}

	dbAdmin := r.entityToModel(admin)

	if err := tx.WithContext(ctx).Save(dbAdmin).Error; err != nil {
		return fmt.Errorf("failed to update admin: %w", err)
	}

	return nil
}

func (r *adminRepository) Delete(ctx context.Context, tx *gorm.DB, admin *entities.Admin) error {
	if admin == nil {
		return errors.New("admin data cannot be empty")
	}

	if admin.User != nil {
		if err := tx.WithContext(ctx).
			Model(&models.User{}).
			Where("id = ?", admin.User.Id).
			Update("is_active", false).Error; err != nil {
			return fmt.Errorf("failed to deactive admin: %w", err)
		}
	}

	if err := tx.WithContext(ctx).
		Model(&models.Admin{}).
		Where("id = ?", admin.Id).
		Update("is_deleted", true).Error; err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}

	return nil
}

func (r *adminRepository) modelToEntity(dbAdmin *models.Admin) *entities.Admin {
	admin := &entities.Admin{
		Id:         dbAdmin.Id,
		UserId:     dbAdmin.UserId,
		AdminName:  dbAdmin.AdminName,
		AdminPhone: dbAdmin.AdminPhone,
		CreatedAt:  dbAdmin.CreatedAt,
		UpdatedAt:  dbAdmin.UpdatedAt,
	}

	if dbAdmin.User != nil {
		admin.User = r.userModelToEntity(dbAdmin.User)
	}

	return admin
}

func (r *adminRepository) entityToModel(admin *entities.Admin) *models.Admin {
	return &models.Admin{
		Id:         admin.Id,
		UserId:     admin.UserId,
		AdminName:  admin.AdminName,
		AdminPhone: admin.AdminPhone,
		CreatedAt:  admin.CreatedAt,
		UpdatedAt:  admin.UpdatedAt,
	}
}

func (r *adminRepository) userModelToEntity(dbUser *models.User) *entities.User {
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
