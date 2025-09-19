package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateObservationAnswersTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE observation_answers (
  			id INTEGER PRIMARY KEY AUTO_INCREMENT,
  			observation_id INTEGER NOT NULL,
  			question_id INTEGER NOT NULL,
  			answer BOOLEAN NOT NULL, 
  			score_earned INTEGER NOT NULL DEFAULT 0, 
  			note TEXT NULL,
  			
  			INDEX observation_aspect_idx (observation_id, question_id),
  			INDEX observation_id_idx (observation_id),
  			INDEX question_id_idx (question_id),
  			
  			UNIQUE KEY unique_observation_aspect (observation_id, question_id),
  			CHECK (score_earned >= 0 AND score_earned <= 3),
  			
  			FOREIGN KEY (observation_id) REFERENCES observations(id) ON DELETE CASCADE,
  			FOREIGN KEY (question_id) REFERENCES observation_questions(id) ON DELETE CASCADE
		);
    `).Error
}

func RollbackCreateObservationAnswersTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE observation_answers;").Error
}
