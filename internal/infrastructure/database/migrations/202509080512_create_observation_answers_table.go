package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateObservationAnswersTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE observation_answers (
			id INTEGER PRIMARY KEY,
			observation_id INTEGER NOT NULL,
			aspect_index VARCHAR(20) NOT NULL,
			answer ENUM('Ya', 'Tidak') NOT NULL,
			note TEXT NULL,
			
			INDEX observation_id_idx (observation_id),
			FOREIGN KEY (observation_id) REFERENCES observations(id)
		);
    `).Error
}

func RollbackCreateObservationAnswersTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE observation_answers;").Error
}
