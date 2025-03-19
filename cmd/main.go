package main

import (
	_ "context"
	_ "database/sql"
	_ "fmt"
	"log"
	_ "net"
	"os"

	grpcserver "github.com/50-Course/notes-tracker/cmd/grpc"
	"github.com/50-Course/notes-tracker/cmd/repository"
	"github.com/50-Course/notes-tracker/scripts/migrations"
	_ "github.com/50-Course/notes-tracker/shared/proto"
	"github.com/50-Course/notes-tracker/shared/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/uptrace/bun"
	_ "github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("Error loading .env file. Please confirm file exists in the right file path and try again.")
	}


	internalServerPort, exists := os.LookupEnv("INTERNAL_SERVER_PORT")
	if !exists {
		log.Printf("INTERNAL_SERVER_PORT not set in environment. Defaulting to 50051")
		// we should just default to some open unused port
		internalServerPort = "50051"
	}

	dbUrl := utils.BuildDatabaseURL()
	db, err := utils.ConnectToDB(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	// run migrations
	log.Println("[Startup] Running database migrations...")
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatalf("[Startup] Migrations failed: %v", err)
	}

	// we would then initialize our grpc server here
	repo := repository.NewTaskRepository(db)
	grpcserver.RunGRPCServer(repo, internalServerPort)
	log.Printf("gRPC Server started on port %s", internalServerPort)
}
