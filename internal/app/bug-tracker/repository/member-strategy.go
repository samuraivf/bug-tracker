package repository

import (
	"database/sql"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
)

//go:generate mockgen -source=member-strategy.go -destination=mocks/member-strategy.go

type member interface {
	IsMember(projectID, userID uint64) error
}

type memberStrategy struct {
	db  *sql.DB
	log log.Log
}

func new_memberStrategy(db *sql.DB, log log.Log) member {
	return &memberStrategy{db, log}
}

func (s *memberStrategy) IsMember(projectID, userID uint64) error {
	result := s.db.QueryRow(
		"SELECT member_id FROM projects_members WHERE project_id = $1 AND member_id = $2",
		projectID,
		userID,
	)
	memberId := 0
	if err := result.Scan(&memberId); err != nil {
		s.log.Error(err)
		return ErrNoRights
	}

	return nil
}
