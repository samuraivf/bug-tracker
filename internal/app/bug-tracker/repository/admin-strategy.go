package repository

import (
	"database/sql"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
)

//go:generate mockgen -source=admin-strategy.go -destination=mocks/admin-strategy.go

type admin interface {
	IsAdmin(projectID, userID uint64) error
}

type adminStrategy struct {
	db  *sql.DB
	log log.Log
}

func new_adminStrategy(db *sql.DB, log log.Log) admin {
	return &adminStrategy{db, log}
}

func (s *adminStrategy) IsAdmin(projectID, userID uint64) error {
	result := s.db.QueryRow("SELECT admin FROM projects WHERE id = $1", projectID)

	var adminID uint64
	if err := result.Scan(&adminID); err != nil {
		s.log.Error(err)
		return err
	}

	if userID != adminID {
		return ErrNoRights
	}

	return nil
}