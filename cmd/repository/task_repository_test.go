package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	_ "time"

	"github.com/50-Course/notes-tracker/shared/models"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func setupTestDB() *bun.DB {
	dsn := os.Getenv("DATABASE_URL")
	// if this is not set we should just default to some test db running locally for now
	if dsn == "" {
		log.Fatalf("Database URL not provided in credentials. Defaulting to local test db")

		dsn = "postgres://postgres:password@localhost:5432/notes_tracker_test?sslmode=disable"
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	testDB := bun.NewDB(sqldb, pgdialect.New())

	if err := testDB.Ping(); err != nil {
		log.Fatalf("Database Connection Error: Cannot connect to database: %v", err)
	}

	// we all love fancy debug logic in debugging/tests
	testDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// apply migrations
	_, _ = testDB.NewDropTable().Model((*models.Task)(nil)).IfExists().IfExists().Cascade().Exec(context.Background())
	_, err := testDB.NewCreateTable().Model((*models.Task)(nil)).IfNotExists().Exec(context.Background())
	if err != nil {
		log.Fatalf("Database Integrity Error: Unable to apply migrations: %v", err)
	}

	// run the tests
	// m.Run()

	// teardown our test data
	// _ = testDB.Close()

	// os.Exit(m.Run())
	return testDB
}

func TestTaskRepository(t *testing.T) {
	testDB := setupTestDB()
	repo := NewTaskRepository(testDB)

	t.Run("Create a new Task Item", func(t *testing.T) {
		task := &models.Task{
			ID:          "test-uuid-4",
			Title:       "Test Task",
			Description: "Although Optional, this is to contain some very long words, :eyes:",
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

	t.Run("Get Task By ID", func(t *testing.T) {
		task := &models.Task{
			Title:       "Fetch a Test Task",
			Description: "This task will be fetched",
		}
		_ = repo.CreateTask(context.Background(), task)

		fetchedTask, err := repo.GetTask(context.Background(), task.ID)
		if err != nil {
			t.Errorf("Single fetch operation failed: %v", err)
		}
		if fetchedTask.ID != task.ID {
			t.Errorf("Expected ID %v, got %v", task.ID, fetchedTask.ID)
		}
		if fetchedTask.Title != task.Title {
			t.Errorf("Expected Title %v, got %v", task.Title, fetchedTask.Title)
		}
	})

	t.Run("List Tasks", func(t *testing.T) {
		tasks, err := repo.ListTasks(context.Background())
		if err != nil {
			t.Errorf("Batch fetch operation failed: %v", err)
		}
		if len(tasks) == 0 {
			t.Errorf("Expected at least one task, but got 0")
		}
	})

	t.Run("Delete a Task", func(t *testing.T) {
		task := &models.Task{
			Title:       "Delete me later",
			Description: "I am a task that will be deleted later",
		}

		_ = repo.CreateTask(context.Background(), task)
		err := repo.DeleteTask(context.Background(), task.ID)
		if err != nil {
			t.Fatalf("Failed to delete task: %v", err)
		}

		_, err = repo.GetTask(context.Background(), "test-uuid-4")
		if err == nil {
			t.Errorf("Expected task to be deleted, but it still exists")
		}
	})

	t.Run("Update Task", func(t *testing.T) {
		task := &models.Task{
			Title:       "Test Task",
			Description: "Commentary before update",
		}

		_ = repo.CreateTask(context.Background(), task)

		task.Title = "Updated Task"
		task.Description = "Post-Update Commentary. Have a great day"

		err := repo.UpdateTask(context.Background(), task)
		if err != nil {
			t.Errorf("Failed to update task: %v", err)
		}

		// recheck if the task was truly updated
		updateTask, err := repo.GetTask(context.Background(), task.ID)
		if err != nil {
			t.Fatalf("Failed to get updated task: %v", err)
		}

		if updateTask.Title != "Updated Task" {
			t.Errorf("Expected task title to be 'Updated Task', got %s", updateTask.Title)
		}

		if updateTask.Description != "Post-Update Commentary. Have a great day" {
			t.Error("Operation Failed: Description was not updated. Expected")
		}
	})
}
