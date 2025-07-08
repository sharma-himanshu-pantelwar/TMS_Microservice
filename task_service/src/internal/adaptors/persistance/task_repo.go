package persistance

import (
	"user_service/src/internal/core/tasks"
)

type TaskRepo struct {
	db *Database
}

func NewTaskRepo(d *Database) tasks.TaskRepoImpl {
	return TaskRepo{db: d}
}

func (t TaskRepo) CreateTask(taskData tasks.TaskDetails) (tasks.TaskDetails, error) {
	var task tasks.TaskDetails
	var taskId int

	query := `
		INSERT INTO TASKS (
			ASSIGNED_BY, ASSIGNED_TO, ASSIGNED_AT, ACCEPTED_AT, DEADLINE, PRIORITY, STATUS
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING ID;
	`

	err := t.db.db.QueryRow(query, taskData.AssignedBy, taskData.AssignedTo, taskData.AssignedAt, taskData.AcceptedAt, taskData.Deadline, taskData.Priority, taskData.Status).Scan(&taskId)
	if err != nil {
		return task, err
	}

	task = taskData
	task.Id = taskId
	return task, nil
}
