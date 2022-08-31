package domain

import (
	"errors"
	"strconv"
	"time"
)

const (
	NoPriority = iota
	LowPriority
	MediumPriority
	HighPriority
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	StartDate   time.Time `json:"startDate"`
	DueDate     time.Time `json:"dueDate"`
	Owner       string    `json:"owner"`
	IsDone      bool      `json:"isDone"`
}

func NewTask(description string, priority int, owner string) Task {
	return Task{
		Description: description,
		Priority:    priority,
		StartDate:   time.Now().UTC(),
		Owner:       owner,
		IsDone:      false,
	}
}

func (t *Task) CreateValidation() error {
	if t.Description == "" {
		return errors.New("task description is empty")
	}
	if t.Priority > HighPriority {
		return errors.New("not valid priority: " + strconv.Itoa(t.Priority))
	}
	if t.Owner == "" {
		return errors.New("owner is empty")
	}

	return nil
}
