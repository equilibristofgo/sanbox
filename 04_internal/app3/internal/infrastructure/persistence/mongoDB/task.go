package mongoDB

import (
	"context"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/domain"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/ports"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type taskMongo struct {
	db *mongo.Client
}

func NewTaskMongoDB(db *mongo.Client) ports.TasksRepository {
	return &taskMongo{db}
}

func (p *taskMongo) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	headerLog := "TaskMongo - Create - "
	log.Println(ctx, headerLog+"Init")

	//TODO: implement

	log.Println(ctx, headerLog+"End")
	return domain.Task{}, nil
}

func (p *taskMongo) Find(ctx context.Context, id int) (domain.Task, error) {
	headerLog := "TaskMongo - Find - "
	log.Println(ctx, headerLog+"Init")

	//TODO: implement

	log.Println(ctx, headerLog+"End")
	return domain.Task{}, nil
}

func (p *taskMongo) Update(ctx context.Context, task domain.Task) error {
	headerLog := "TaskMongo - Update - "
	log.Println(ctx, headerLog+"Init")

	//TODO: implement

	log.Println(ctx, headerLog+"End")
	return nil
}

func (p *taskMongo) Delete(ctx context.Context, id int) error {
	headerLog := "TaskMongo - Delete - "
	log.Println(ctx, headerLog+"Init")

	//TODO: implement

	log.Println(ctx, headerLog+"End")
	return nil
}
