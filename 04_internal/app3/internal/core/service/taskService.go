package service

import (
	"context"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/domain"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/ports"
	"time"
)

type taskService struct {
	taskRepository ports.TasksRepository
}

func NewTaskService(taskRepository ports.TasksRepository) ports.TaskService {
	return &taskService{taskRepository}
}

func (s *taskService) Get(ctx context.Context, id int) (domain.Task, error) {
	return s.taskRepository.Find(ctx, id)
}

func (s *taskService) Create(ctx context.Context, description string, priority int, owner string) (domain.Task, error) {
	task := domain.Task{
		Description: description,
		Priority:    priority,
		Owner:       owner,
	}
	return s.taskRepository.Create(ctx, task)
}

func (s *taskService) Update(ctx context.Context, id int, done bool) error {
	task := domain.Task{
		Id:      id,
		DueDate: time.Now().UTC(),
		IsDone:  done,
	}
	return s.taskRepository.Update(ctx, task)
}

func (s *taskService) Delete(ctx context.Context, id int) error {
	return s.taskRepository.Delete(ctx, id)
}
