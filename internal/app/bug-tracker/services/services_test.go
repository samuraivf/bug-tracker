package services

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mock_redis "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
	mock_repository "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository/mocks"
)

func Test_NewService(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	auth := NewAuth()
	repo := &repository.Repository{User: mock_repository.NewMockUser(c), Project: mock_repository.NewMockProject(c)}
	redis := mock_redis.NewMockRedis(c)

	expected := &Service{Auth: auth, Redis: NewRedis(redis), User: NewUser(repo.User), Project: NewProject(repo.Project)}

	require.Equal(t, expected, NewService(repo, redis))
}
