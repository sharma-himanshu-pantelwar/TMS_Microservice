package tasks

import "time"

type TaskDetails struct {
	Id         int       `json:"id"`
	AssignedBy int       `json:"assigned_by"`
	AssignedTo int       `json:"assigned_to,omitempty"`
	AssignedAt time.Time `json:"assigned_at"`
	AcceptedAt time.Time `json:"accepted_at"`
	Deadline   time.Time `json:"deadline"`
	Priority   string    `json:"priority,omitempty"`
	Status     string    `json:"status,omitempty"`
}

type TaskRepoImpl interface {
	CreateTask(taskData TaskDetails) (TaskDetails, error)
	UpdateTask(taskData TaskDetails, taskId int) (TaskDetails, error)
}
type TaskServiceImpl interface {
	CreateTask(taskData TaskDetails) (TaskDetails, error)
	UpdateTask(taskData TaskDetails, taskId int) (TaskDetails, error)
}
