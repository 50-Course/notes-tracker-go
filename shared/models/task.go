package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Represents a Todo Item; effectively a task or an action
// triaged to be done by us later
type Task struct {
	bun.BaseModel `bun:"table:tasks,alias:t"`

	ID          string `bun:",pk,type:uuid,default:gen_random_uuid()"`
	Title       string `bun:",notnull"`
	Description string
	CreatedAt   time.Time `bun:",default:current_timestamp"`
	UpdatedAt   bun.NullTime
}

// formats to pretty representation
func (t *Task) String() string {
	return t.Title
}
