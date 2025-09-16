package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateTherapistsTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE therapists(
			id      			CHAR(26) PRIMARY KEY NOT NULL,
			user_id 			CHAR(26)             NOT NULL,
			therapist_name    	VARCHAR(100)         NOT NULL,
			therapist_section 	ENUM('Okupasi', 'Fisio', 'Wicara', 'Paedagog') NOT NULL,
			therapist_phone 	VARBINARY(100)       NOT NULL,
			created_at 			TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
			updated_at 			TIMESTAMP            DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			is_deleted BOOLEAN DEFAULT FALSE,
			CONSTRAINT fk_therapists_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
			UNIQUE KEY  idx_therapists_user_id (user_id)
		);
    `).Error
}

func RollbackCreateTherapistsTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE therapists;").Error
}
