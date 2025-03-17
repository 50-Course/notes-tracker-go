package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/50-Course/notes-tracker/shared/models"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var testDB *bun.DB

func testMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	// if this is not set we should just default to some test db running locally for now
	if dsn == "" {
		log.Fatalf("Database URL not provided in credentials. Defaulting to local test db")

		dsn = "postgres://postgres:password@localhost:5432/notes_tracker_test?sslmode=disable"
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	testDB := bun.NewDB(sqldb, pgdialect.New())

	// we all love fancy debug logic in debugging/tests
	testDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// apply migrations
	_, _ = testDB.NewDropTable().Model((*models.Task)(nil)).IfExists().IfExists().Cascade().Exec(context.Background())
	_, _ = testDB.NewCreateTable().Model((*models.Task)(nil)).IfNotExists().Exec(context.Background())

	// run the tests
	m.Run()

	// teardown
	_ = testDB.Close()

	os.Exit(m.Run())
}

func TestTaskRepository(t *testing.T) {
	repo := NewTaskRepository(testDB)

	t.Run("Create a new Task Item", func(t *testing.T) {
		task := &models.Task{
			ID:          "test-uuid-4",
			Title:       "Test Task",
			Description: "Although Optional, this is to contain some very long words, :eyes:",
			CreatedAt:   time.Now(),
		}

		err := repo.CreateTask(context.Background(), task)
		if err != nil {
			t.Errorf("Error creating task: %v", err)
		}

		// check if the task ID is set
		if task.ID == "" {
			t.Errorf("Expected task ID to be set, got 0")
		}

		// task must exist in our database
		savedTask, err := repo.GetTask(context.Background(), "test-uuid-4")
		if err != nil {
			t.Fatalf("Task does not exist. Error: %v", err)
		}

		if savedTask.Title != task.Title {
			t.Errorf("Expected task title to be %s, got %s", task.Title, savedTask.Title)
		}
	})

	t.Run("Delete a Task", func(t *testing.T) {
		err := repo.DeleteTask(context.Background(), "test-uuid-4")
		if err != nil {
			t.Fatalf("Failed to delete task: %v", err)
		}

		_, err = repo.GetTask(context.Background(), "test-uuid-4")
		if err == nil {
			t.Errorf("Expected task to be deleted, but it still exists")
		}
	})
}
