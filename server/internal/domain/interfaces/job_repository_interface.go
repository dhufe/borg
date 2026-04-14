package interfaces

import (
	"context"
	"lath/borg/internal/domain/model"
)

type JobRepository interface {
	FindById(ctx context.Context, id uint) (*model.Job, error)
	Create(ctx context.Context, job *model.Job) error
	Update(ctx context.Context, job *model.Job) error
	Delete(ctx context.Context, id uint) error

	FindAll(ctx context.Context) ([]model.Job, error)
	FindPendingTasks(ctx context.Context, tasks *[]model.Job) error
	FindTasksForCleanUp(ctx context.Context) ([]model.Job, error)

	Count(ctx context.Context) (int, error)
	CountByStatus(ctx context.Context, status model.JobStatus) (int, error)
	CountByAllStatuses(ctx context.Context) (map[model.JobStatus]int, error)
}
