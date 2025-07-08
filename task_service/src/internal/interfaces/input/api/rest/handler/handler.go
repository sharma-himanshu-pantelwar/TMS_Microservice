package handler

import (
	"encoding/json"
	"net/http"
	"user_service/src/internal/core/tasks"
	"user_service/src/pkg/response"
)

type TaskHandler struct {
	taskService tasks.TaskServiceImpl
}

func NewTaskHandler(taskService tasks.TaskServiceImpl) TaskHandler {
	return TaskHandler{
		taskService: taskService,
	}
}

func (th TaskHandler) RegisterTaskHandler(w http.ResponseWriter, r *http.Request) {
	var taskData tasks.TaskDetails
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusBadRequest,
			Error:          err.Error(),
			Message:        "failed to create user",
		}
		response.Set()
		return
	}

	insertedTask, err := th.taskService.CreateTask(taskData)
	if err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusInternalServerError,
			Error:          err.Error(),
		}
		response.Set()
		return
	}
	response := response.Response{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		Message:        "User created successfully",
		Error:          "none",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Data: insertedTask,
	}
	response.Set()

}
