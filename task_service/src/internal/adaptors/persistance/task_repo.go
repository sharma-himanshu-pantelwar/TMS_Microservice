package persistance

import (
	"fmt"
	"task_service/src/internal/core/tasks"
	"time"
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
			ASSIGNED_BY, ASSIGNED_TO, ASSIGNED_AT,TASK_NAME,TASK_DESCRIPTION, DEADLINE, PRIORITY, STATUS
		) VALUES ($1, $2, $3, $4, $5, $6, $7,$8)
		RETURNING ID;
	`

	fmt.Printf("CreateTask - taskData: %+v\n", taskData)
	err := t.db.db.QueryRow(query, taskData.AssignedBy, taskData.AssignedTo, taskData.AssignedAt, taskData.TaskName, taskData.TaskDescription, taskData.Deadline, taskData.Priority, taskData.Status).Scan(&taskId)
	if err != nil {
		fmt.Printf("CreateTask - Error: %v\n", err)
		return task, err
	}

	task = taskData
	task.Id = taskId
	fmt.Printf("CreateTask - Created task with ID: %d\n", taskId)
	return task, nil
}

func (t TaskRepo) GetAllTasks(userId int64) ([]tasks.TaskDetails, error) {
	var allTasks []tasks.TaskDetails

	query := `
		SELECT ID, ASSIGNED_BY, ASSIGNED_TO, ASSIGNED_AT, TASK_NAME, TASK_DESCRIPTION, DEADLINE, PRIORITY, STATUS
		FROM TASKS
		WHERE (ASSIGNED_BY = $1 OR ASSIGNED_TO = $1) AND IS_TRASH='FALSE'; 
	`

	// fmt.Printf("GetAllTasks - userId: %d, query: %s\n", userId, query)
	rows, err := t.db.db.Query(query, userId)
	if err != nil {
		// fmt.Printf("GetAllTasks - Query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task tasks.TaskDetails
		err := rows.Scan(
			&task.Id,
			&task.AssignedBy,
			&task.AssignedTo,
			&task.AssignedAt,
			&task.TaskName,
			&task.TaskDescription,
			&task.Deadline,
			&task.Priority,
			&task.Status,
		)
		if err != nil {
			// fmt.Printf("GetAllTasks - Scan error: %v\n", err)
			return nil, err
		}
		// fmt.Printf("GetAllTasks - Found task: %+v\n", task)
		allTasks = append(allTasks, task)
	}

	if err = rows.Err(); err != nil {
		// fmt.Printf("GetAllTasks - Rows error: %v\n", err)
		return nil, err
	}

	fmt.Printf("GetAllTasks - Total tasks found: %d\n", len(allTasks))
	return allTasks, nil
}
func (t TaskRepo) GetMyTasks(userId int64) ([]tasks.TaskDetails, error) {
	var allTasks []tasks.TaskDetails

	query := `
		SELECT ID, ASSIGNED_BY, ASSIGNED_TO, ASSIGNED_AT, TASK_NAME, TASK_DESCRIPTION, DEADLINE, PRIORITY, STATUS,IS_TRASH
		FROM TASKS
		WHERE ASSIGNED_TO = $1; 
	`
	// AND IS_TRASH=FALSE

	// fmt.Printf("GetAllTasks - userId: %d, query: %s\n", userId, query)
	rows, err := t.db.db.Query(query, userId)
	if err != nil {
		// fmt.Printf("GetMyTask - Query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task tasks.TaskDetails
		err := rows.Scan(
			&task.Id,
			&task.AssignedBy,
			&task.AssignedTo,
			&task.AssignedAt,
			&task.TaskName,
			&task.TaskDescription,
			&task.Deadline,
			&task.Priority,
			&task.Status,
			&task.IsTrash,
		)
		if err != nil {
			// fmt.Printf("GetMyTasks - Scan error: %v\n", err)
			return nil, err
		}
		// fmt.Printf("GetMyTasks - Found task: %+v\n", task)
		allTasks = append(allTasks, task)
	}

	if err = rows.Err(); err != nil {
		// fmt.Printf("GetMyTasks - Rows error: %v\n", err)
		return nil, err
	}

	// fmt.Printf("GetMyTasks - Total tasks found: %d\n", len(allTasks))
	return allTasks, nil
}

func (t TaskRepo) UpdateTask(taskData tasks.TaskDetails, taskId int) (tasks.TaskDetails, error) {
	var existing tasks.TaskDetails

	// Fetch the existing task
	selectQuery := `
		SELECT ID, ASSIGNED_BY, ASSIGNED_TO, TASK_NAME, TASK_DESCRIPTION, ASSIGNED_AT, DEADLINE, PRIORITY, STATUS,IS_TRASH
		FROM TASKS WHERE ID = $1;
	`
	err := t.db.db.QueryRow(selectQuery, taskId).Scan(
		&existing.Id,
		&existing.AssignedBy,
		&existing.AssignedTo,
		&existing.TaskName,
		&existing.TaskDescription,
		&existing.AssignedAt,
		&existing.Deadline,
		&existing.Priority,
		&existing.Status,
		&existing.IsTrash,
	)
	if err != nil {
		return existing, err
	}

	// Update  non-empty values
	// if taskData.AssignedBy != 0 {
	// 	existing.AssignedBy = taskData.AssignedBy
	// }
	if taskData.AssignedTo != 0 {
		existing.AssignedTo = taskData.AssignedTo
	}
	if !taskData.AssignedAt.IsZero() {
		existing.AssignedAt = taskData.AssignedAt
	}
	if !taskData.Deadline.IsZero() {
		existing.Deadline = taskData.Deadline
	}
	if taskData.Priority != "" {
		existing.Priority = taskData.Priority
	}
	if taskData.Status != "" {
		existing.Status = taskData.Status
	}
	if taskData.TaskName != "" {
		existing.TaskName = taskData.TaskName
	}
	if taskData.TaskDescription != "" {
		existing.TaskDescription = taskData.TaskDescription
	}

	// Update
	updateQuery := `
		UPDATE TASKS SET
			ASSIGNED_BY = $1,
			ASSIGNED_TO = $2,
			ASSIGNED_AT = $3,
			DEADLINE = $4,
			PRIORITY = $5,
			STATUS = $6,
			TASK_NAME = $7,
			TASK_DESCRIPTION = $8
		WHERE ID = $9
		RETURNING ID, ASSIGNED_BY, ASSIGNED_TO, TASK_NAME, TASK_DESCRIPTION, ASSIGNED_AT, DEADLINE, PRIORITY, STATUS, IS_TRASH;
	`

	var updated tasks.TaskDetails
	err = t.db.db.QueryRow(
		updateQuery,
		existing.AssignedBy,
		existing.AssignedTo,
		existing.AssignedAt,
		existing.Deadline,
		existing.Priority,
		existing.Status,
		existing.TaskName,
		existing.TaskDescription,
		taskId,
	).Scan(
		&updated.Id,
		&updated.AssignedBy,
		&updated.AssignedTo,
		&updated.TaskName,
		&updated.TaskDescription,
		&updated.AssignedAt,
		&updated.Deadline,
		&updated.Priority,
		&updated.Status,
		&updated.IsTrash,
	)
	if err != nil {
		// fmt.Println("Get my tasks error before return ", err)
		return updated, err
	}
	return updated, nil
}

func (t TaskRepo) DeleteTask(userId int64, taskId int) (tasks.TaskDetails, error) {

	// Update
	deleteQuery := `
		UPDATE TASKS SET
			IS_TRASH = TRUE
		WHERE ID = $1 AND ASSIGNED_BY=$2
		RETURNING ID, ASSIGNED_BY, ASSIGNED_TO, TASK_NAME, TASK_DESCRIPTION, ASSIGNED_AT, DEADLINE, PRIORITY, STATUS;
	`
	var deleted tasks.TaskDetails
	err := t.db.db.QueryRow(
		deleteQuery,
		taskId,
		userId,
	).Scan(
		&deleted.Id,
		&deleted.AssignedBy,
		&deleted.AssignedTo,
		&deleted.TaskName,
		&deleted.TaskDescription,
		&deleted.AssignedAt,
		&deleted.Deadline,
		&deleted.Priority,
		&deleted.Status,
	)
	if err != nil {
		return deleted, err
	}
	return deleted, nil
}

func (t TaskRepo) GetAllTasksInBin(userId int64) ([]tasks.TaskDetails, error) {
	var allBinTasks []tasks.TaskDetails
	query := `
		SELECT ID, ASSIGNED_BY, ASSIGNED_TO, ASSIGNED_AT, TASK_NAME, TASK_DESCRIPTION, DEADLINE, PRIORITY, STATUS, IS_TRASH
		FROM TASKS
		WHERE (ASSIGNED_BY = $1 OR ASSIGNED_TO = $1) AND IS_TRASH='TRUE'; 
	`
	// fmt.Printf("GetAllTasks - userId: %d, query: %s\n", userId, query)
	rows, err := t.db.db.Query(query, userId)
	if err != nil {
		// fmt.Printf("GetAllTasks - Query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task tasks.TaskDetails
		err := rows.Scan(
			&task.Id,
			&task.AssignedBy,
			&task.AssignedTo,
			&task.AssignedAt,
			&task.TaskName,
			&task.TaskDescription,
			&task.Deadline,
			&task.Priority,
			&task.Status,
			&task.IsTrash,
		)
		if err != nil {
			// fmt.Printf("GetAllTasks - Scan error: %v\n", err)
			return nil, err
		}
		// fmt.Printf("GetAllTasks - Found task: %+v\n", task)
		allBinTasks = append(allBinTasks, task)
	}

	if err = rows.Err(); err != nil {
		// fmt.Printf("GetAllTasks - Rows error: %v\n", err)
		return nil, err
	}
	fmt.Printf("GetAllTasksInBin - Total tasks found: %d\n", len(allBinTasks))
	return allBinTasks, nil
}

func (t TaskRepo) RestoreTaskFromBin(userId int64, taskId int) (tasks.TaskDetails, error) {
	// Update
	restoreQuery := `
		UPDATE TASKS SET
			IS_TRASH = FALSE
		WHERE ID = $1 AND ASSIGNED_BY=$2
		RETURNING ID, ASSIGNED_BY, ASSIGNED_TO, TASK_NAME, TASK_DESCRIPTION, ASSIGNED_AT, DEADLINE, PRIORITY, STATUS;
	`
	var restored tasks.TaskDetails
	err := t.db.db.QueryRow(
		restoreQuery,
		taskId,
		userId,
	).Scan(
		&restored.Id,
		&restored.AssignedBy,
		&restored.AssignedTo,
		&restored.TaskName,
		&restored.TaskDescription,
		&restored.AssignedAt,
		&restored.Deadline,
		&restored.Priority,
		&restored.Status,
	)
	if err != nil {
		return restored, err
	}

	return restored, nil
}
func (t TaskRepo) DeleteTaskFromBin(userId int64, taskId int) (tasks.TaskDetails, error) {

	// Update
	deleteQuery := `
		DELETE FROM TASKS 
		WHERE ID = $1 AND ASSIGNED_BY=$2 AND IS_TRASH=TRUE
		RETURNING ID, ASSIGNED_BY, ASSIGNED_TO, TASK_NAME, TASK_DESCRIPTION, ASSIGNED_AT, DEADLINE, PRIORITY, STATUS;
	`

	var delFromBin tasks.TaskDetails
	err := t.db.db.QueryRow(
		deleteQuery,
		taskId,
		userId,
	).Scan(
		&delFromBin.Id,
		&delFromBin.AssignedBy,
		&delFromBin.AssignedTo,
		&delFromBin.TaskName,
		&delFromBin.TaskDescription,
		&delFromBin.AssignedAt,
		&delFromBin.Deadline,
		&delFromBin.Priority,
		&delFromBin.Status,
	)
	if err != nil {
		return delFromBin, err
	}

	return delFromBin, nil
}
func (t TaskRepo) DeleteTaskPermanently(userId int64, taskId int) (tasks.TaskDetails, error) {

	// Update
	deleteQuery := `
		DELETE FROM TASKS 
		WHERE ID = $1 AND ASSIGNED_BY=$2 AND IS_TRASH=FALSE
		RETURNING ID, ASSIGNED_BY, ASSIGNED_TO, TASK_NAME, TASK_DESCRIPTION, ASSIGNED_AT, DEADLINE, PRIORITY, STATUS;
	`

	var delFrom tasks.TaskDetails
	err := t.db.db.QueryRow(
		deleteQuery,
		taskId,
		userId,
	).Scan(
		&delFrom.Id,
		&delFrom.AssignedBy,
		&delFrom.AssignedTo,
		&delFrom.TaskName,
		&delFrom.TaskDescription,
		&delFrom.AssignedAt,
		&delFrom.Deadline,
		&delFrom.Priority,
		&delFrom.Status,
	)
	if err != nil {
		return delFrom, err
	}

	return delFrom, nil
}
func (t TaskRepo) CheckAssignedUserStatus(userId int64, deadline time.Time) (bool, int, error) {
	var count int

	query := `SELECT COUNT(*) FROM TASKS WHERE ASSIGNED_TO=$1 AND DEADLINE::DATE=$2::DATE;`

	err := t.db.db.QueryRow(query, userId, deadline).Scan(&count)
	if err != nil {
		return false, 0, err
	}
	if count > 3 {
		return false, count, nil
	}
	return true, count, nil
}
