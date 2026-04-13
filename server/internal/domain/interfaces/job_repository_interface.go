package interfaces

import (
	"context"
	"lath/borg/internal/domain/model"
)

type JobRepository interface {
	FindById(ctx context.Context, id uint) (model.Job, error)
	Create(ctx context.Context, job *model.Job) error
	Update(ctx context.Context, job *model.Job) error
	Delete(ctx context.Context, id uint) error
}
