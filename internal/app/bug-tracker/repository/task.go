package repository

import (
	"database/sql"
	"time"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
)

type TaskRepository struct {
	db  *sql.DB
	log log.Log
}

func NewTaskRepo(db *sql.DB, log log.Log) Task {
	return &TaskRepository{
		db:  db,
		log: log,
	}
}

func (r *TaskRepository) CreateTask(taskData *dto.CreateTaskDto) (uint64, error) {
	result := r.db.QueryRow(
		"INSERT INTO tasks (name, description, task_priority, project_id, task_type, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		taskData.Name,
		taskData.Description,
		taskData.TaskPriority,
		taskData.ProjectID,
		taskData.TaskType,
		time.Now(),
	)

	var taskID uint64
	if err := result.Scan(&taskID); err != nil {
		r.log.Error(err)
		return 0, err
	}
	r.log.Infof("Create task: id = %d", taskID)

	return taskID, nil
}
