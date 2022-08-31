package ports

import (
	"context"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/domain"
)

type TasksRepository interface {
	Create(ctx context.Context, task domain.Task) (domain.Task, error)
	Find(ctx context.Context, id int) (domain.Task, error)
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id int) error
}

type TaskService interface {
	Get(ctx context.Context, id int) (domain.Task, error)
	Create(ctx context.Context, description string, priority int, owner string) (domain.Task, error)
	Update(ctx context.Context, id int, done bool) error
	Delete(ctx context.Context, id int) error
}
