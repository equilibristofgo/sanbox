package controller

import (
	"context"
	"encoding/json"
	"github.com/equilibristofgo/sandbox/04_internal/app3/apps/http/controller/dto"
	"github.com/equilibristofgo/sandbox/04_internal/app3/internal/core/ports"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type TaskHttp struct {
	taskService ports.TaskService
}

func NewTaskHttp(router *chi.Mux, taskServer ports.TaskService) {
	r := &TaskHttp{taskServer}
	router.Post("/", r.createTask)
	router.Get("/{taskId}", r.getTask)
	router.Put("/{taskId}", r.updateTask)
	router.Delete("/{taskId}", r.deleteTask)
}

// createTask godoc
// @Summary New task
// @Description Create new task
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param task body dto.Task true "task body"
// @Success 200 {object} dto.Task
// @Router /v1/template/tasks [post]
func (c *TaskHttp) createTask(w http.ResponseWriter, r *http.Request) {
	headerLog := "TaskHttp - createTask - "
	w.Header().Set("Content-Type", "application/json")
	log.Println(headerLog + "Init")

	var task dto.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Println(headerLog + "Error: Invalid input, object invalid")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid input, object invalid")
		return
	}

	response, err := c.taskService.Create(context.TODO(), task.Description, task.Priority, task.Owner)
	if err != nil {
		log.Println(headerLog + "Error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	log.Println(headerLog + "End")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// getTask godoc
// @Summary Get task
// @Description Get task with specified id
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param taskId path string true "taskId"
// @Success 200 {object} dto.Task
// @Router /v1/template/tasks/{taskId} [get]
func (c *TaskHttp) getTask(w http.ResponseWriter, r *http.Request) {
	headerLog := "TaskHttp - getTask - "
	w.Header().Set("Content-Type", "application/json")
	log.Println(headerLog + "Init")

	taskIdS := chi.URLParam(r, "taskId")
	taskId, _ := strconv.Atoi(taskIdS)

	response, err := c.taskService.Get(context.TODO(), taskId)
	if err != nil {
		log.Println(headerLog + "Error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	log.Println(headerLog + "End")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// updateTask godoc
// @Summary Update task
// @Description Update task with specified id
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param taskId path string true "taskId"
// @Param task body dto.Task true "task body"
// @Success 200 {object} string
// @Router /v1/template/tasks/{taskId} [put]
func (c *TaskHttp) updateTask(w http.ResponseWriter, r *http.Request) {
	headerLog := "TaskHttp - updateTask - "
	w.Header().Set("Content-Type", "application/json")
	log.Println(headerLog + "Init")

	taskIdS := chi.URLParam(r, "taskId")
	taskId, _ := strconv.Atoi(taskIdS)

	var task dto.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Println(headerLog + "Error: Invalid input, object invalid")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid input, object invalid")
		return
	}

	err = c.taskService.Update(context.TODO(), taskId, task.IsDone)
	if err != nil {
		log.Println(headerLog + "Error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	log.Println(headerLog + "End")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("OK")
}

// deleteTask godoc
// @Summary Delete task
// @Description Delete task with specified id
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param taskId path string true "taskId"
// @Success 200 {object} string
// @Router /v1/template/tasks/{taskId} [delete]
func (c *TaskHttp) deleteTask(w http.ResponseWriter, r *http.Request) {
	headerLog := "TaskHttp - deleteTask - "
	w.Header().Set("Content-Type", "application/json")
	log.Println(headerLog + "Init")

	taskIdS := chi.URLParam(r, "taskId")
	taskId, _ := strconv.Atoi(taskIdS)

	err := c.taskService.Delete(context.TODO(), taskId)
	if err != nil {
		log.Println(headerLog + "Error: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	log.Println(headerLog + "End")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Task " + taskIdS + " has been deleted")
}
