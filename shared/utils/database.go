package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

// Constructs a databse connection string from environment variables.
func BuildDatabaseURL() string {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")

	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("Missing required database environment variables")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
}

// Initializes a database connection using Bun ORM
func ConnectToDB(dbUrl string) (*bun.DB, error) {
	log.Println("[DB] Connecting to database with URL:", dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("[DB] Error connecting to database: %v", err)
	}

	// Ping DB to confirm connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[DB] Failed to connect to database: %v", err)
	}

	log.Println("[DB] Connected to database successfully")
	return bun.NewDB(db, pgdialect.New()), nil
}
