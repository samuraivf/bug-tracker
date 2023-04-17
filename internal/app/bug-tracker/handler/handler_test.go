package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/golang/mock/gomock"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
	mock_kafka "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/kafka/mocks"
	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
	mock_services "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services/mocks"
)

func Test_newHandler(t *testing.T) {
	handler := NewHandler(nil, nil, nil)
	require.Equal(t, &Handler{nil, nil, nil}, handler)

	c := gomock.NewController(t)
	auth := mock_services.NewMockAuth(c)
	srv := &services.Service{Auth: auth}
	log := mock_log.NewMockLog(c)
	kafka := mock_kafka.NewMockKafka(c)
	handler = NewHandler(srv, log, kafka)

	require.Equal(t, &Handler{
		service: &services.Service{Auth: auth},
		log: log,
		kafka: kafka,
	}, handler)
}