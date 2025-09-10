package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateVerificationCodesTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE verification_codes (
			id         INTEGER  PRIMARY KEY NOT NULL AUTO_INCREMENT,
			user_id    CHAR(26)             NOT NULL,
			code       VARCHAR(6)           NOT NULL,
			status 	   ENUM('Pending', 'Used', 'Revoked') DEFAULT 'Pending' NOT NULL ,
			expires_at DATETIME             NOT NULL,
			created_at TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_verification_codes_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
			INDEX idx_verification_codes_user_id (user_id)
		);
    `).Error
}

func RollbackCreateVerificationCodesTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE verification_codes;").Error
}
