package postgreSQL

import (
	"context"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/domain"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/ports"
	"gorm.io/gorm"
	"log"
)

type taskPostgres struct {
	db *gorm.DB
}

func NewTaskPostgres(db *gorm.DB) ports.TasksRepository {
	return &taskPostgres{db}
}

func (p *taskPostgres) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	headerLog := "TaskPostgres - Create - "
	log.Println(headerLog + "Init")

	//TODO: implement

	log.Println(headerLog + "End")
	return domain.Task{}, nil
}

func (p *taskPostgres) Find(ctx context.Context, id int) (domain.Task, error) {
	headerLog := "TaskPostgres - Find - "
	log.Println(ctx, headerLog+"Init")

	//TODO: implement

	log.Println(ctx, headerLog+"End")
	return domain.Task{}, nil
}

func (p *taskPostgres) Update(ctx context.Context, task domain.Task) error {
	headerLog := "TaskPostgres - Update - "
	log.Println(ctx, headerLog+"Init")

	//TODO: implement

	log.Println(ctx, headerLog+"End")
	return nil
}

func (p *taskPostgres) Delete(ctx context.Context, id int) error {
	headerLog := "TaskPostgres - Delete - "
	log.Println(ctx, headerLog+"Init")

	//TODO: implement

	log.Println(ctx, headerLog+"End")
	return nil
}
