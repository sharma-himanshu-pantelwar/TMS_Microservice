package usecase

import (
	"context"
	"errors"
	"fmt"
	redisadaptor "task_service/src/internal/adaptors/redis"
	"task_service/src/internal/core/tasks"
)

type TaskService struct {
	taskRepo       tasks.TaskRepoImpl
	redisPublisher *redisadaptor.RedisPublisher
}

func NewTaskService(taskRepo tasks.TaskRepoImpl, redisPublisher *redisadaptor.RedisPublisher) TaskService {
	return TaskService{
		taskRepo:       taskRepo,
		redisPublisher: redisPublisher,
	}
}

func (ts TaskService) CreateTask(taskData tasks.TaskDetails) (tasks.TaskDetails, error) {
	createdTask, err := ts.taskRepo.CreateTask(taskData)
	if err != nil {
		// fmt.Println(err)
		return createdTask, errors.New("failed to create task, try again later")
	}

	//publish notification
	if ts.redisPublisher != nil {
		msg := fmt.Sprintf("Task Created: %+v", createdTask)
		_ = ts.redisPublisher.PublishTaskNotification(context.Background(), "task_notifications", msg)
		fmt.Println("message sent to redis : ", msg)
	}

	return createdTask, nil
}
func (ts TaskService) GetAllTasks(userId int64) ([]tasks.TaskDetails, error) {
	allTasks, err := ts.taskRepo.GetAllTasks(userId)
	if err != nil {
		// fmt.Println(err)
		return allTasks, errors.New("failed to fetch all task, try again later")
	}

	return allTasks, nil
}
func (ts TaskService) GetMyTasks(userId int64) ([]tasks.TaskDetails, error) {
	allTasks, err := ts.taskRepo.GetMyTasks(userId)
	if err != nil {
		// fmt.Println(err)
		return allTasks, errors.New("failed to fetch all task, try again later")
	}

	return allTasks, nil
}

func (ts TaskService) UpdateTask(taskData tasks.TaskDetails, taskId int) (tasks.TaskDetails, error) {
	updatedTask, err := ts.taskRepo.UpdateTask(taskData, taskId)
	if err != nil {
		fmt.Println(err)
		return updatedTask, errors.New("failed to update task, try again later")
	}
	//publish notification
	if ts.redisPublisher != nil {
		msg := fmt.Sprintf("Updated Task: %+v", updatedTask)
		_ = ts.redisPublisher.PublishTaskNotification(context.Background(), "task_notifications", msg)
		fmt.Println("message sent to redis : ", msg)
	}
	return updatedTask, nil
}
func (ts TaskService) DeleteTask(userId int64, taskId int) (tasks.TaskDetails, error) {
	updatedTask, err := ts.taskRepo.DeleteTask(userId, taskId)
	if err != nil {
		fmt.Println(err)
		return updatedTask, errors.New("failed to update task, try again later")
	}
	if ts.redisPublisher != nil {
		msg := fmt.Sprintf("Deleted Task: %+v", updatedTask)
		_ = ts.redisPublisher.PublishTaskNotification(context.Background(), "task_notifications", msg)
		fmt.Println("message sent to redis : ", msg)
	}
	return updatedTask, nil
}

// GetAllTasksInBin(userId)
func (ts TaskService) GetAllTasksInBin(userId int64) ([]tasks.TaskDetails, error) {
	allTasks, err := ts.taskRepo.GetAllTasksInBin(userId)
	if err != nil {
		// fmt.Println(err)
		return allTasks, errors.New("failed to fetch all task, try again later")
	}

	return allTasks, nil
}
func (ts TaskService) RestoreTask(userId int64, taskId int) (tasks.TaskDetails, error) {
	allTasks, err := ts.taskRepo.RestoreTaskFromBin(userId, taskId)
	if err != nil {
		// fmt.Println(err)
		return allTasks, errors.New("failed to fetch all task, try again later")
	}

	return allTasks, nil
}

func (ts TaskService) DeleteTaskFromBin(userId int64, taskId int) (tasks.TaskDetails, error) {
	updatedTask, err := ts.taskRepo.DeleteTaskFromBin(userId, taskId)
	if err != nil {
		fmt.Println(err)
		return updatedTask, errors.New("failed to update task, try again later")
	}

	return updatedTask, nil
}
func (ts TaskService) DeleteTaskPermanently(userId int64, taskId int) (tasks.TaskDetails, error) {
	updatedTask, err := ts.taskRepo.DeleteTaskPermanently(userId, taskId)
	if err != nil {
		fmt.Println(err)
		return updatedTask, errors.New("failed to update task, try again later")
	}

	return updatedTask, nil
}
