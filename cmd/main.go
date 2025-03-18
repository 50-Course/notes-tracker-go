package main

import (
	"database/sql"
	"fmt"
	_ "fmt"
	"log"
	_ "net"
	"os"

	_ "github.com/lib/pq"
	grpcserver "github.com/50-Course/notes-tracker/cmd/grpc"
	"github.com/50-Course/notes-tracker/cmd/repository"
	_ "github.com/50-Course/notes-tracker/shared/proto"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func connectToDB(dbUrl string) (*bun.DB, error) {
	db, err := sql.Open("pg", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %v", err)
	}
	return bun.NewDB(db, pgdialect.New()), nil
}

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("Error loading .env file. Please confirm file exists in the right file path and try again.")
	}

	dbUrl, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		log.Fatal("DATABASE_URL not set in environment")
	}

	internalServerPort, exists := os.LookupEnv("INTERNAL_SERVER_PORT")
	if !exists {
		log.Printf("INTERNAL_SERVER_PORT not set in environment. Defaulting to 50051")
		// we should just default to some open unused port
		internalServerPort = "50051"
	}

	db, err := connectToDB(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	// we would then initialize our grpc server here
	repo := repository.NewTaskRepository(db)
	grpcserver.RunGRPCServer(repo, internalServerPort)
	log.Printf("gRPC Server started on port %s", internalServerPort)
}
