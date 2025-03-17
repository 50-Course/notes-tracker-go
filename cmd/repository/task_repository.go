package repository

import (
	"context"

	"github.com/50-Course/notes-tracker/shared/models"
	"github.com/uptrace/bun"
)

type TaskRepository struct {
	db *bun.DB
}

func NewTaskRepository(db *bun.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	_, err := r.db.NewInsert().Model(task).Exec(ctx)
	return err
}

func (r *TaskRepository) GetTask(ctx context.Context, id string) (*models.Task, error) {
	task := new(models.Task)
	err := r.db.NewSelect().Model(task).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// func (r *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
// 	_, err := r.db.NewUpdate().Model(task).Where("id = ?", task.ID).Exec(ctx)
// 	return err
// }

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model((*models.Task)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
