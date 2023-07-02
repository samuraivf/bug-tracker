package repository

import (
	"database/sql"
	"time"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
)

type TaskRepository struct {
	db    *sql.DB
	log   log.Log
	admin admin
}

func NewTaskRepo(db *sql.DB, log log.Log, admin admin) Task {
	return &TaskRepository{
		db:    db,
		log:   log,
		admin: admin,
	}
}

func (r *TaskRepository) CreateTask(taskData *dto.CreateTaskDto, userID uint64) (uint64, error) {
	if err := r.admin.IsAdmin(taskData.ProjectID, userID); err != nil {
		return 0, err
	}

	result := r.db.QueryRow(
		"INSERT INTO tasks (name, description, task_priority, project_id, task_type, created_at, perform_to) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		taskData.Name,
		taskData.Description,
		taskData.TaskPriority,
		taskData.ProjectID,
		taskData.TaskType,
		time.Now(),
		taskData.PerformTo,
	)

	var taskID uint64
	if err := result.Scan(&taskID); err != nil {
		r.log.Error(err)
		return 0, err
	}
	r.log.Infof("Create task: id = %d", taskID)

	return taskID, nil
}

func (r *TaskRepository) WorkOnTask(workOnTaskData *dto.WorkOnTaskDto, userID uint64) error {
	result := r.db.QueryRow(
		"SELECT member_id FROM projects_members WHERE project_id = $1 AND member_id = $2",
		workOnTaskData.ProjectID,
		userID,
	)
	memberId := 0
	if result.Scan(&memberId) != nil && r.admin.IsAdmin(workOnTaskData.ProjectID, userID) != nil {
		return ErrNoRights
	}

	_, err := r.db.Exec("UPDATE tasks SET assignee = $1 WHERE id = $2 AND assignee IS NULL", userID, workOnTaskData.TaskID)
	if err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

func (r *TaskRepository) StopWorkOnTask(workOnTaskData *dto.WorkOnTaskDto, userID uint64) error {
	result := r.db.QueryRow(
		"SELECT member_id FROM projects_members WHERE project_id = $1 AND member_id = $2",
		workOnTaskData.ProjectID,
		userID,
	)
	memberId := 0
	if result.Scan(&memberId) != nil && r.admin.IsAdmin(workOnTaskData.ProjectID, userID) != nil {
		return ErrNoRights
	}

	_, err := r.db.Exec("UPDATE tasks SET assignee = NULL WHERE id = $1 AND assignee IS NOT NULL", workOnTaskData.TaskID)
	if err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

func (r *TaskRepository) UpdateTask(taskData *dto.UpdateTaskDto, userID uint64) (uint64, error) {
	if err := r.admin.IsAdmin(taskData.ProjectID, userID); err != nil {
		return 0, err
	}

	result := r.db.QueryRow(
		"UPDATE tasks SET name = $1, description = $2, task_priority = $3, project_id = $4, task_type = $5, perform_to = $6 WHERE id = $7 RETURNING id",
		taskData.Name,
		taskData.Description,
		taskData.TaskPriority,
		taskData.ProjectID,
		taskData.TaskType,
		taskData.PerformTo,
		taskData.TaskID,
	)

	var taskID uint64
	if err := result.Scan(&taskID); err != nil {
		r.log.Error(err)
		return 0, err
	}
	r.log.Infof("Create task: id = %d", taskID)

	return taskID, nil
}

func (r *TaskRepository) GetTaskById(id uint64) (*models.Task, error) {
	result := r.db.QueryRow(
		`SELECT 
			id, 
			name, 
			description, 
			task_priority, 
			project_id, 
			task_type, 
			assignee, 
			created_at, 
			perform_to 
		FROM tasks WHERE id = $1`,
		id,
	)

	task := new(models.Task)
	if err := result.Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Priority,
		&task.ProjectID,
		&task.TaskType,
		&task.Assignee,
		&task.CreatedAt,
		&task.PerformTo,
	); err != nil {
		r.log.Error(err)
		return nil, err
	}
	r.log.Infof("Get task: id = %d", id)

	return task, nil
}

func (r *TaskRepository) DeleteTask(taskData *dto.DeleteTaskDto, userID uint64) error {
	if err := r.admin.IsAdmin(taskData.ProjectID, userID); err != nil {
		return err
	}

	_, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", taskData.TaskID)

	if err != nil {
		r.log.Error(err)
	}

	return err
}
