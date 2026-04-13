package persistence

import (
	"context"
	"lath/borg/internal/domain/interfaces"
	"lath/borg/internal/domain/model"

	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
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
