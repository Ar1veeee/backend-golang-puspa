package database

import (
	"backend-golang/shared/config"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
	dbName := config.GetEnv("DB_NAME", "")
	if dbName == "" {
		log.Fatal().Msg("DB_NAME environment variable not set")
	}

	dbUser := config.GetEnv("DB_USER", "root")
	dbPass := config.GetEnv("DB_PASS", "")
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "3306")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	configLogger := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		DisableForeignKeyConstraintWhenMigrating: false,
	}

	db, err := gorm.Open(mysql.Open(dsn), configLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info().Msgf("Database dialect configured: %s", db.Dialector.Name())

	log.Info().Msg("Database connected successfully!")
	return db, nil
}
