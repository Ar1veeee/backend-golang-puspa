package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateParentDetailsTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE parent_details (
			id                			CHAR(26) PRIMARY KEY 					NOT NULL,
			parent_id         			CHAR(26)             					NOT NULL,
			parent_type       			ENUM('Ayah', 'Ibu', 'Wali')  			NOT NULL,
			parent_name       			VARCHAR(100)         					NOT NULL,
			parent_phone      			VARBINARY(100)       					NOT NULL,
			parent_age 		  			INTEGER 							    NULL,
			parent_occupation 			VARCHAR(100)                          	NULL,
			relationship_with_child		VARCHAR(100) 							NULL, 
			CONSTRAINT fk_parent_details_parent FOREIGN KEY (parent_id) REFERENCES parents (id) ON DELETE CASCADE,
			INDEX             idx_parent_details_parent_id (parent_id),
			INDEX             idx_parent_details_parent_type (parent_type)
		);
    `).Error
}

func RollbackCreateParentDetailsTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE parent_details;").Error
}
