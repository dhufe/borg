package interfaces

import (
	"context"
	"lath/borg/internal/domain/model"
)

type UserRepository interface {
	FindById(ctx context.Context, id uint) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)

	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}
