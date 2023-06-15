package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
)

func Test_NewRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	c := gomock.NewController(t)
	defer c.Finish()

	log := mock_log.NewMockLog(c)
	admin := new_adminStrategy(db, log)
	expectedRepo := &Repository{
		User:    NewUserRepo(db, log),
		Project: NewProjectRepo(db, log, admin),
		Task:    NewTaskRepo(db, log, admin),
	}
	repo := NewRepository(db, log)

	require.Equal(t, expectedRepo, repo)
}
