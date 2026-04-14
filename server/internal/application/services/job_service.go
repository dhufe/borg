package services

import (
	"context"

	"lath/borg/internal/domain/interfaces"
	"lath/borg/internal/domain/model"
)

type JobService struct {
	repo            interfaces.JobRepository
	fileStoragePath string
}

func NewTaskService(repo interfaces.JobRepository, fileStoragePath string) *JobService {
	return &JobService{
		repo:            repo,
		fileStoragePath: fileStoragePath,
	}
}

func (s *JobService) CreateTask(
	ctx context.Context,
	filename string) (*model.Job, error) {
	task := &model.Job{
		FileName: filename,
		Status:   model.StatusPending,
	}

	err := s.repo.Create(ctx, task)
	return task, err
}

func (s *JobService) FileStoragePath() string {
	return s.fileStoragePath
}

func (s *JobService) FindPendingTasks(ctx context.Context) ([]model.Job, error) {
	var tasks []model.Job
	err := s.repo.FindPendingTasks(ctx, &tasks)
	return tasks, err
}

func (s *JobService) GetTask(ctx context.Context, taskID uint) (*model.Job, error) {
	return s.repo.FindById(ctx, taskID)
}

func (s *JobService) DeleteTask(ctx context.Context, taskID uint) error {
	return s.repo.Delete(ctx, taskID)
}

func (s *JobService) GetTasksForCleanUp(ctx context.Context) ([]model.Job, error) {
	tasks, err := s.repo.FindTasksForCleanUp(ctx)
	return tasks, err
}

func (s *JobService) GetAllJobs(ctx context.Context) ([]model.Job, error) {
	tasks, err := s.repo.FindAll(ctx)
	return tasks, err
}

func (s *JobService) UpdateTask(
	ctx context.Context,
	taskId uint,
	filename string,
	status model.JobStatus,
) (*model.Job, error) {
	task := &model.Job{
		ID:       taskId,
		FileName: filename,
		Status:   status,
	}

	err := s.repo.Update(ctx, task)
	return task, err
}
