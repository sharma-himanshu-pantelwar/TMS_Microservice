package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task_service/src/internal/core/tasks"
	samw "task_service/src/internal/interfaces/input/api/rest/middleware"
	"task_service/src/pkg/response"

	"github.com/go-chi/chi/v5"
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
	userIdVal := r.Context().Value(samw.UserIdKey)
	// fmt.Println("userId val from context is ", userIdVal)
	if userIdVal == nil {
		http.Error(w, "userId not found in context ", http.StatusInternalServerError)
		return
	}
	var taskData tasks.TaskDetails

	// Cast userIdVal to int64 since ValidateSession returns int64
	taskData.AssignedBy = userIdVal.(int64)
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusBadRequest,
			Error:          err.Error(),
			Message:        "failed to create task",
		}
		response.Set()
		return
	}
	//
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

func (th TaskHandler) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	userIdVal := r.Context().Value(samw.UserIdKey)
	fmt.Println("userId val from context is ", userIdVal)
	if userIdVal == nil {
		http.Error(w, "userId not found in context ", http.StatusInternalServerError)
		return
	}

	userId := userIdVal.(int64)
	fmt.Printf("Converted userId to int64: %d\n", userId)
	//
	allTasks, err := th.taskService.GetAllTasks(userId)
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
		Message:        "Tasks fetched successfully",
		Error:          "none",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Data: allTasks,
	}
	response.Set()

}

func (th TaskHandler) GetMyTasksHandler(w http.ResponseWriter, r *http.Request) {
	userIdVal := r.Context().Value(samw.UserIdKey)
	fmt.Println("userId val from context is ", userIdVal)
	if userIdVal == nil {
		http.Error(w, "userId not found in context ", http.StatusInternalServerError)
		return
	}

	userId := userIdVal.(int64)
	fmt.Printf("Converted userId to int64: %d\n", userId)
	//
	allTasks, err := th.taskService.GetMyTasks(userId)
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
		Message:        "Tasks fetched successfully",
		Error:          "none",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Data: allTasks,
	}
	response.Set()

}

func (th TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var updatedTaskData tasks.TaskDetails
	taskIdString := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(taskIdString)
	if err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusBadRequest,
			Error:          err.Error(),
			Message:        "invalid taskid",
		}
		response.Set()
		return
	}
	err = json.NewDecoder(r.Body).Decode(&updatedTaskData)
	if err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusBadRequest,
			Error:          err.Error(),
			Message:        "invalid req body",
		}
		response.Set()
		return
	}

	// Get userID from context and set AssignedBy
	userIdVal := r.Context().Value(samw.UserIdKey)
	if userIdVal == nil {
		http.Error(w, "userId not found in context", http.StatusUnauthorized)
		return
	}
	updatedTaskData.AssignedBy = userIdVal.(int64)

	updatedTask, err := th.taskService.UpdateTask(updatedTaskData, taskId)
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
		Data: updatedTask,
	}
	response.Set()

}
func (th TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskIdString := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(taskIdString)
	if err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusBadRequest,
			Error:          err.Error(),
			Message:        "invalid taskid",
		}
		response.Set()
		return
	}
	// err = json.NewDecoder(r.Body).Decode(&updatedTaskData)
	// if err != nil {
	// 	response := response.Response{
	// 		ResponseWriter: w,
	// 		StatusCode:     http.StatusBadRequest,
	// 		Error:          err.Error(),
	// 		Message:        "invalid req body",
	// 	}
	// 	response.Set()
	// 	return
	// }

	// Get userID from context and set AssignedBy
	userIdVal := r.Context().Value(samw.UserIdKey)
	if userIdVal == nil {
		http.Error(w, "userId not found in context", http.StatusUnauthorized)
		return
	}
	// updatedTaskData.AssignedBy = userIdVal.(int64)

	deletedTask, err := th.taskService.DeleteTask(userIdVal.(int64), taskId)
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
		Data: deletedTask,
	}
	response.Set()

}

func (th TaskHandler) GetTaskBinHandler(w http.ResponseWriter, r *http.Request) {

	userIdVal := r.Context().Value(samw.UserIdKey)
	fmt.Println("userId val from context is ", userIdVal)
	if userIdVal == nil {
		http.Error(w, "userId not found in context ", http.StatusInternalServerError)
		return
	}

	userId := userIdVal.(int64)
	fmt.Printf("Converted userId to int64: %d\n", userId)
	//
	allTasksInBin, err := th.taskService.GetAllTasksInBin(userId)
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
		Message:        "Tasks fetched successfully from bin",
		Error:          "none",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Data: allTasksInBin,
	}
	response.Set()

}

func (th TaskHandler) RestoreTaskFromBinHandler(w http.ResponseWriter, r *http.Request) {
	taskIdString := chi.URLParam(r, "id")
	taskId, err := strconv.Atoi(taskIdString)
	if err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusBadRequest,
			Error:          err.Error(),
			Message:        "invalid taskid",
		}
		response.Set()
		return
	}

	userIdVal := r.Context().Value(samw.UserIdKey)
	if userIdVal == nil {
		http.Error(w, "userId not found in context", http.StatusUnauthorized)
		return
	}
	// updatedTaskData.AssignedBy = userIdVal.(int64)

	deletedTask, err := th.taskService.RestoreTask(userIdVal.(int64), taskId)
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
		Message:        "Task restored from bin successfully",
		Error:          "none",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Data: deletedTask,
	}
	response.Set()

}
