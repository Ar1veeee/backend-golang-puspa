package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateObservationsTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE observations (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			child_id CHAR(26) NOT NULL,
			therapist_id CHAR(26) NULL,
			scheduled_date DATE NOT NULL,
			age_category ENUM('Balita', 'Anak-anak', 'Remaja', 'Lainnya') NOT NULL,
			total_score INTEGER NULL,
			conclusion TEXT NULL,
			recommendation TEXT NULL,
			status ENUM('Pending', 'Scheduled', 'Complete') NOT NULL DEFAULT 'Pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

			INDEX child_id_status_idx (child_id, status),
			FOREIGN KEY (child_id) REFERENCES childrens(id)
		);
    `).Error
}

func RollbackCreateObservationsTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE observations;").Error
}
