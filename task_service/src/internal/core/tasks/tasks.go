package tasks

import "time"

type TaskDetails struct {
	Id              int       `json:"id"`
	AssignedBy      int64     `json:"assigned_by"`
	AssignedTo      int64     `json:"assigned_to,omitempty"`
	AssignedAt      time.Time `json:"assigned_at"`
	TaskName        string    `json:"task_name"`
	TaskDescription string    `json:"task_description"`
	// AcceptedAt      time.Time `json:"accepted_at"`
	Deadline time.Time `json:"deadline"`
	Priority string    `json:"priority,omitempty"`
	Status   string    `json:"status,omitempty"`
	IsTrash  bool      `json:"isTrash"`
}

type TaskRepoImpl interface {
	CreateTask(taskData TaskDetails) (TaskDetails, error)
	UpdateTask(taskData TaskDetails, taskId int) (TaskDetails, error)
	GetAllTasks(userId int64) ([]TaskDetails, error)
	GetMyTasks(userId int64) ([]TaskDetails, error)
	DeleteTask(userId int64, taskId int) (TaskDetails, error)
	GetAllTasksInBin(userId int64) ([]TaskDetails, error)
	RestoreTaskFromBin(userId int64, taskId int) (TaskDetails, error)
	DeleteTaskFromBin(userId int64, taskId int) (TaskDetails, error)
	DeleteTaskPermanently(userId int64, taskId int) (TaskDetails, error)
	CheckAssignedUserStatus(userId int64, deadline time.Time) (bool, int, error)
}
type TaskServiceImpl interface {
	CreateTask(taskData TaskDetails) (TaskDetails, error)
	UpdateTask(taskData TaskDetails, taskId int) (TaskDetails, error)
	GetAllTasks(userId int64) ([]TaskDetails, error)
	GetMyTasks(userId int64) ([]TaskDetails, error)
	DeleteTask(userId int64, taskId int) (TaskDetails, error)
	DeleteTaskFromBin(userId int64, taskId int) (TaskDetails, error)
	DeleteTaskPermanently(userId int64, taskId int) (TaskDetails, error)
	GetAllTasksInBin(userId int64) ([]TaskDetails, error)
	RestoreTask(userId int64, taskId int) (TaskDetails, error)
	CheckAssignedUserStatus(userId int64, deadline time.Time) (bool, int, error)
}
