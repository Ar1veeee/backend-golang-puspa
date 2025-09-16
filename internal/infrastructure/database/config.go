package database

import (
	"backend-golang/internal/infrastructure/config"
	"errors"
	"time"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	Timezone        string
}

func NewConfig() *Config {
	return &Config{
		Host:            config.GetEnv("DB_HOST", "localhost"),
		Port:            config.GetEnv("DB_PORT", "3306"),
		User:            config.GetEnv("DB_USER", "root"),
		Password:        config.GetEnv("DB_PASS", ""),
		Name:            config.GetEnv("DB_NAME", ""),
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		Timezone:        "Asia/Jakarta",
	}
}

func (c *Config) Validate() error {
	if c.Name == "" {
		return errors.New("database name is required")
	}
	if c.User == "" {
		return errors.New("database user is required")
	}
	if c.Host == "" {
		return errors.New("database host is required")
	}
	return nil
}
