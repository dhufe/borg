package repositories

import (
	"context"
	"lath/borg/internal/domain/interfaces"
	"lath/borg/internal/domain/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (u UserRepository) FindById(ctx context.Context, id uint) (model.Job, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) FindByEmail(ctx context.Context, email string) (model.Job, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) Create(ctx context.Context, job *model.Job) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) Update(ctx context.Context, job *model.Job) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) Delete(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}
