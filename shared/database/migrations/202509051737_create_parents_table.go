package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateParentsTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE parents (
			id      				CHAR(26) PRIMARY KEY NOT NULL,
			user_id 				CHAR(26)             NULL,
			temp_email  			VARCHAR(100)		 NOT NULL,
			registration_status		ENUM('Pending', 'Complete') DEFAULT 'Pending' NOT NULL,
			created_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			
			CONSTRAINT fk_parents_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
			UNIQUE KEY idx_user_id (user_id),
			UNIQUE KEY idx_temp_email (temp_email)
		);
    `).Error
}

func RollbackCreateParentsTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE parents;").Error
}
