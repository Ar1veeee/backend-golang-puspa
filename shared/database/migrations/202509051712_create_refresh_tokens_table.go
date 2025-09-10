package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateRefreshTokensTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE refresh_tokens (
			id         INTEGER  PRIMARY KEY NOT NULL AUTO_INCREMENT,
			user_id    CHAR(26)             NOT NULL,
			token      TEXT                 NOT NULL,
			expires_at DATETIME             NOT NULL,
			created_at TIMESTAMP            NOT NULL DEFAULT CURRENT_TIMESTAMP,
			revoked    BOOLEAN              NOT NULL DEFAULT FALSE,
			CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
			INDEX idx_refresh_tokens_user_id (user_id)
		);
    `).Error
}

func RollbackCreateRefreshTokensTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE refresh_tokens;").Error
}
