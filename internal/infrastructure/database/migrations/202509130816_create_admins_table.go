package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateAdminsTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE admins(
			id      			CHAR(26) PRIMARY KEY NOT NULL,
			user_id 			CHAR(26)             NOT NULL,
			admin_name    	VARCHAR(100)         NOT NULL,
			admin_phone 	VARBINARY(100)       NOT NULL,
			created_at 			TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
			updated_at 			TIMESTAMP            DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			is_deleted BOOLEAN DEFAULT FALSE,
			CONSTRAINT fk_admins_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
			UNIQUE KEY  idx_admins_user_id (user_id)
		);
    `).Error
}

func RollbackCreateAdminsTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE admins;").Error
}
