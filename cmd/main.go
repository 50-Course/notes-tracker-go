package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "net"
	"os"

	grpcserver "github.com/50-Course/notes-tracker/cmd/grpc"
	"github.com/50-Course/notes-tracker/cmd/repository"
	_ "github.com/50-Course/notes-tracker/shared/proto"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

func buildDatabaseURL() string {
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

func connectToDB(dbUrl string) (*bun.DB, error) {
	log.Println("Connecting to database with URL:", dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %v", err)
	}

	// ping db to truely know if we are connecting to the databse
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")
	return bun.NewDB(db, pgdialect.New()), nil
}

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("Error loading .env file. Please confirm file exists in the right file path and try again.")
	}

	// dbUrl, exists := os.LookupEnv("DATABASE_URL")
	// if !exists {
	// 	log.Fatal("DATABASE_URL not set in environment")
	// }

	internalServerPort, exists := os.LookupEnv("INTERNAL_SERVER_PORT")
	if !exists {
		log.Printf("INTERNAL_SERVER_PORT not set in environment. Defaulting to 50051")
		// we should just default to some open unused port
		internalServerPort = "50051"
	}

	dbUrl := buildDatabaseURL()
	db, err := connectToDB(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	// we would then initialize our grpc server here
	repo := repository.NewTaskRepository(db)
	grpcserver.RunGRPCServer(repo, internalServerPort)
	log.Printf("gRPC Server started on port %s", internalServerPort)
}
