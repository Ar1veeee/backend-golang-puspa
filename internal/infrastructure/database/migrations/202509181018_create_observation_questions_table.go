package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateObservationQuestionsTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE observation_questions (
		  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  		  question_code VARCHAR(6) UNIQUE NOT NULL,
		  age_category ENUM('Balita', 'Anak-anak', 'Remaja', 'Lainya') NOT NULL,
		  question_number INTEGER NOT NULL, 
		  question_text TEXT NOT NULL, 
		  score INTEGER NOT NULL, 
		  is_active BOOLEAN DEFAULT TRUE,
		  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		  
		  INDEX age_category_number_idx (age_category, question_number),
		  INDEX question_code_idx (question_code),
		  
		  UNIQUE KEY unique_age_question_number (age_category, question_number),
		  UNIQUE KEY unique_question_code (question_code),
		  CHECK (score >= 1 AND score <= 3)
		);
    `).Error
}

func RollbackCreateObservationQuestionsTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE observation_questions;").Error
}
