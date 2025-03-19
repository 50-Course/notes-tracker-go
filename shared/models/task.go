package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Represents a Todo Item; effectively a task or an action
// triaged to be done by us later
//
// @Description: Task Model
type Task struct {
	bun.BaseModel `bun:"table:tasks,alias:t" swaggerignore:"true"`

	ID          string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Title       string `bun:",notnull"`
	Description string
	CreatedAt   time.Time    `bun:",default:current_timestamp"`
	UpdatedAt   bun.NullTime `swaggertype:"string" format:"date-time"`
}

// formats to pretty representation
func (t *Task) String() string {
	return t.Title
}

// does the pretty formating of our timestamps
func (t *Task) AsTime(timestamp string) time.Time {
	tm, _ := time.Parse(time.RFC3339, timestamp)
	return tm
}

func (t *Task) GetCreatedAt() time.Time {
	return t.AsTime(t.CreatedAt.String())
}

func (t *Task) GetUpdatedAt() time.Time {
	return t.AsTime(t.UpdatedAt.String())
}

// Defines the request payload for creating a task.
type TaskRequest struct {
	Title       string `json:"title" example:"Buy groceries"`
	Description string `json:"description" example:"Milk, Bread, Eggs"`
}

// Defines the response payload for returning a task.
type TaskResponse struct {
	ID          string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title       string `json:"title" example:"Buy groceries"`
	Description string `json:"description" example:"Milk, Bread, Eggs"`
	CreatedAt   string `json:"created_at" example:"2025-03-19T08:58:10.605Z"`
	UpdatedAt   string `json:"updated_at" example:"2025-03-19T08:58:10.605Z"`
}
