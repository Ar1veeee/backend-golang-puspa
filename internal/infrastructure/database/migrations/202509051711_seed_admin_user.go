package migrations

import (
	helpers2 "backend-golang/internal/helpers"
	"backend-golang/pkg/models"
	"time"

	"gorm.io/gorm"
)

func SeedUsersTableUp(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&models.User{}).Where("username = ?", "admin1").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hashedPassword, err := helpers2.HashPassword("Admin_Pusp4.")
	if err != nil {
		return err
	}

	admin := models.User{
		Id:        helpers2.GenerateULID(),
		Username:  "admin1",
		Email:     "admin@gmail.com",
		Password:  string(hashedPassword),
		Role:      "Admin",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return tx.Create(&admin).Error
}

func SeedUsersTableDown(tx *gorm.DB) error {
	return tx.Where("username = ?", "admin1").Delete(&models.User{}).Error
}
