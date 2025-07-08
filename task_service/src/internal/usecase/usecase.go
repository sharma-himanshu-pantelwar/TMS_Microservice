package usecase

import (
	"errors"
	"user_service/src/internal/core/tasks"
)

type TaskService struct {
	taskRepo tasks.TaskRepoImpl
}

func NewTaskService(taskRepo tasks.TaskRepoImpl) tasks.TaskRepoImpl {
	return TaskService{taskRepo: taskRepo}
}

func (ts TaskService) CreateTask(taskData tasks.TaskDetails) (tasks.TaskDetails, error) {
	createdTask, err := ts.taskRepo.CreateTask(taskData)
	if err != nil {
		// fmt.Println(err)
		return createdTask, errors.New("failed to create task, try again later")
	}

	return createdTask, nil
}

func (ts TaskService) UpdateTask(taskData tasks.TaskDetails, taskId int) (tasks.TaskDetails, error) {
	updatedTask, err := ts.taskRepo.UpdateTask(taskData, taskId)
	if err != nil {
		// fmt.Println(err)
		return updatedTask, errors.New("failed to update task, try again later")
	}

	return updatedTask, nil
}
