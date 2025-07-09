package persistance

import (
	"task_service/src/internal/core/tasks"
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

func (t TaskRepo) UpdateTask(taskData tasks.TaskDetails, taskId int) (tasks.TaskDetails, error) {
	var task tasks.TaskDetails

	query := `
		UPDATE TASKS SET
			ASSIGNED_BY = $1,
			ASSIGNED_TO = $2,
			ASSIGNED_AT = $3,
			ACCEPTED_AT = $4,
			DEADLINE = $5,
			PRIORITY = $6,
			STATUS = $7
		WHERE ID = $8
		RETURNING ID, ASSIGNED_BY, ASSIGNED_TO, ASSIGNED_AT, ACCEPTED_AT, DEADLINE, PRIORITY, STATUS;
	`

	err := t.db.db.QueryRow(
		query,
		taskData.AssignedBy,
		taskData.AssignedTo,
		taskData.AssignedAt,
		taskData.AcceptedAt,
		taskData.Deadline,
		taskData.Priority,
		taskData.Status,
		taskId,
	).Scan(
		&task.Id,
		&task.AssignedBy,
		&task.AssignedTo,
		&task.AssignedAt,
		&task.AcceptedAt,
		&task.Deadline,
		&task.Priority,
		&task.Status,
	)
	if err != nil {
		return task, err
	}

	return task, nil
}
