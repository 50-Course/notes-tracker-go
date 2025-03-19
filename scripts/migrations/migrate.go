package migrations

import (
	"context"
	"log"

	"github.com/50-Course/notes-tracker/shared/models"
	_ "github.com/50-Course/notes-tracker/shared/utils"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	_ "github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

// Heavily inspired by Django's makemigrations command
// This function will create the migrations for the models or tables
// that are not yet migrated
func RunMigrations(db *bun.DB) error {
	ctx := context.Background()

	// here we keep trrrack of the models to migrate
	installedSchemas := []interface{}{
		// Add your models here
		(*models.Task)(nil),
	}

	for _, model := range installedSchemas {
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			// log.Panicf("[Migrations] Failed to apply migrations: %v", err)
			return err
		}

		log.Printf("[Migrations] Migrated model: %T", model)
	}

	log.Println("[Migrations] Migrations completed successfully")
	return nil
}

