package database

import (
	"backend-golang/internal/infrastructure/database/migrations"
	"fmt"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) RunMigrations() error {
	dialect := m.db.Dialector.Name()
	if dialect != "mysql" {
		return fmt.Errorf("unsupported SQL dialect: %s, expected mysql", dialect)
	}

	migrator := gormigrate.New(m.db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID:       "202509051710_create_users_table",
			Migrate:  migrations.MigrateCreateUsersTable,
			Rollback: migrations.RollbackCreateUsersTable,
		},
		{
			ID:       "202509051711_seed_admin_user",
			Migrate:  migrations.SeedUsersTableUp,
			Rollback: migrations.SeedUsersTableDown,
		},
		{
			ID:       "202509051712_create_refresh_tokens_table",
			Migrate:  migrations.MigrateCreateRefreshTokensTable,
			Rollback: migrations.RollbackCreateRefreshTokensTable,
		},
		{
			ID:       "202509051737_create_parents_table",
			Migrate:  migrations.MigrateCreateParentsTable,
			Rollback: migrations.RollbackCreateParentsTable,
		},
		{
			ID:       "202509051739_create_parent_details_table",
			Migrate:  migrations.MigrateCreateParentDetailsTable,
			Rollback: migrations.RollbackCreateParentDetailsTable,
		},
		{
			ID:       "202509051742_create_childrens_table",
			Migrate:  migrations.MigrateCreateChildrensTable,
			Rollback: migrations.RollbackCreateChildrensTable,
		},
		{
			ID:       "202509051744_create_therapists_table",
			Migrate:  migrations.MigrateCreateTherapistsTable,
			Rollback: migrations.RollbackCreateTherapistsTable,
		},
		{
			ID:       "202509071113_create_verification_codes_table",
			Migrate:  migrations.MigrateCreateVerificationCodesTable,
			Rollback: migrations.RollbackCreateVerificationCodesTable,
		},
		{
			ID:       "202509080509_create_observations_table",
			Migrate:  migrations.MigrateCreateObservationsTable,
			Rollback: migrations.RollbackCreateObservationsTable,
		},
		{
			ID:       "202509080512_create_observation_answers_table",
			Migrate:  migrations.MigrateCreateObservationAnswersTable,
			Rollback: migrations.RollbackCreateObservationAnswersTable,
		},
		{
			ID:       "202509130816_create_admins_table",
			Migrate:  migrations.MigrateCreateAdminsTable,
			Rollback: migrations.RollbackCreateAdminsTable,
		},
	})

	if err := migrator.Migrate(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migrated successfully")
	return nil
}
