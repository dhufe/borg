package repositories

import (
	"context"
	"lath/borg/internal/domain/interfaces"
	"lath/borg/internal/domain/model"

	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func (j JobRepository) FindPendingTasks(ctx context.Context, tasks *[]model.Job) error {
	//TODO implement me
	panic("implement me")
}

func (j JobRepository) FindTasksForCleanUp(ctx context.Context) ([]model.Job, error) {
	//TODO implement me
	panic("implement me")
}

func (j JobRepository) Count(ctx context.Context) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (j JobRepository) CountByStatus(ctx context.Context, status model.JobStatus) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (j JobRepository) CountByAllStatuses(ctx context.Context) (map[model.JobStatus]int, error) {
	//TODO implement me
	panic("implement me")
}

func (j JobRepository) FindById(ctx context.Context, id uint) (model.Job, error) {
	//TODO implement me
	panic("implement me")
}

func (j JobRepository) Create(ctx context.Context, job *model.Job) error {
	//TODO implement me
	panic("implement me")
}

func (j JobRepository) Update(ctx context.Context, job *model.Job) error {
	//TODO implement me
	panic("implement me")
}

func (j JobRepository) Delete(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}

func NewJobRepository(db *gorm.DB) interfaces.JobRepository {
	return &JobRepository{db: db}
}
