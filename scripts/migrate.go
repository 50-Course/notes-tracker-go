package main

import (
	"context"
	"log"

	"github.com/50-Course/notes-tracker/shared/models"
	"github.com/50-Course/notes-tracker/shared/utils"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
)

// Heavily inspired by Django's makemigrations command
// This function will create the migrations for the models or tables
// that are not yet migrated
func runMigrations(db *bun.DB) {
	ctx := context.Background()

	// here we keep trrrack of the models to migrate
	installedSchemas := []interface{}{
		// Add your models here
		(*models.Task)(nil),
	}

	for _, model := range installedSchemas {
		if _, err := db.NewCreateTable().Model(model).Exec(ctx); err != nil {
			log.Panicf("[Migrations] Failed to apply migrations: %v", err)
		}

		log.Printf("[Migrations] Migrated model: %T", model)
	}

	log.Println("[Migrations] Migrations completed successfully")
}

// func migrate(db *bun.DB) {
// 	ctx := context.Background()
//
// 	log.Println("[Migrations] Running migration...")
// 	_, err := db.NewCreateTable().Model((*models.Task)(nil)).IfNotExists().Exec(ctx)
// 	if err != nil {
// 		log.Fatalf("[Migrations] Migration failed: %v", err)
// 	}
//
// 	log.Println("[Migrations] Migration successful!")
// }

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := utils.BuildDatabaseURL()
	db, err := utils.ConnectToDB(dbUrl)
	if err != nil {
		// just rethrow the error
		log.Fatal(err)
	}
	runMigrations(db)
}
