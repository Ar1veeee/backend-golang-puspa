package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connection interface {
	GetDB() *gorm.DB
	Close() error
	Ping(ctx context.Context) error
}

type connection struct {
	db     *gorm.DB
	sqlDB  *sql.DB
	config *Config
}

func NewConnection(cfg *Config) (Connection, error) {
	conn := &connection{config: cfg}
	if err := conn.connect(); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *connection) connect() error {
	dsn := c.buildDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: c.setupLogger(),
		NowFunc: func() time.Time {
			loc, _ := time.LoadLocation(c.config.Timezone)
			return time.Now().In(loc)
		},
	})
	if err != nil {
		return err
	}

	c.db = db
	c.sqlDB, _ = db.DB()
	c.configurePool()

	return nil
}

func (c *connection) buildDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		c.config.User, c.config.Password, c.config.Host, c.config.Port,
		c.config.Name, url.QueryEscape(c.config.Timezone))
}

func (c *connection) setupLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

func (c *connection) configurePool() {
	if c.sqlDB == nil {
		return
	}

	c.sqlDB.SetMaxOpenConns(25)
	c.sqlDB.SetMaxIdleConns(25)
	c.sqlDB.SetConnMaxLifetime(5 * time.Minute)
}

func (c *connection) GetDB() *gorm.DB {
	return c.db
}

func (c *connection) Close() error {
	if c.sqlDB != nil {
		return c.sqlDB.Close()
	}
	return nil
}

func (c *connection) Ping(ctx context.Context) error {
	if c.sqlDB != nil {
		return c.sqlDB.PingContext(ctx)
	}
	return fmt.Errorf("sqlDB is not initialized")
}
