package migrations

import (
	"gorm.io/gorm"
)

func MigrateCreateChildrensTable(tx *gorm.DB) error {
	return tx.Exec(`
        CREATE TABLE childrens (
			id                		CHAR(26) PRIMARY KEY NOT NULL,
			parent_id         		CHAR(26) 			 NOT NULL,
			child_name        		VARCHAR(100)         NOT NULL,
			child_gender            ENUM('Laki-laki', 'Perempuan') NOT NULL,
			child_birth_place 		VARCHAR(100)         NOT NULL,
			child_birth_date  		DATE                 NOT NULL,
			child_address           VARBINARY(500)       NOT NULL,
			child_complaint         VARCHAR(200)         NOT NULL,
			child_school            VARCHAR(100) 		 NULL,
			child_service_choice    VARCHAR(250)         NOT NULL,
			child_religion    		ENUM('Islam','Kristen','Katolik','Hindu','Budha','Konghucu','Lainnya') NULL ,
			created_at TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP                                    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_children_parent FOREIGN KEY (parent_id) REFERENCES parents (id) ON DELETE CASCADE,
			INDEX idx_children_parent_id (parent_id)
		);
    `).Error
}

func RollbackCreateChildrensTable(tx *gorm.DB) error {
	return tx.Exec("DROP TABLE childrens;").Error
}
