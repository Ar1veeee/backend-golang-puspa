package migrations

import "gorm.io/gorm"

func MigrateCreateUsersTable(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE users (
			id         CHAR(26) PRIMARY KEY,
			username   VARCHAR(50)                         NOT NULL,
			email      VARCHAR(100)                        NOT NULL,
			password   VARCHAR(255)                        NOT NULL,
			role       ENUM ('Admin', 'Terapis', 'User')   NOT NULL DEFAULT 'User',
			is_active  BOOL                                NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			last_login TIMESTAMP 						   NULL,
			UNIQUE KEY idx_users_username (username),
			UNIQUE KEY idx_users_email (email),
			INDEX idx_users_is_active (is_active)
		);
	`).Error
}

func RollbackCreateUsersTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE users;").Error
}
