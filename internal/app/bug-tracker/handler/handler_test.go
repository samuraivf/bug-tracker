package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mock_kafka "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/kafka/mocks"
	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
	mock_services "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services/mocks"
)

func Test_newHandler(t *testing.T) {
	handler := NewHandler(nil, nil, nil, nil)
	require.Equal(t, &Handler{nil, nil, nil, nil}, handler)

	c := gomock.NewController(t)
	auth := mock_services.NewMockAuth(c)
	srv := &services.Service{Auth: auth}
	log := mock_log.NewMockLog(c)
	kafka := mock_kafka.NewMockKafka(c)
	p := &params{}
	handler = NewHandler(srv, log, kafka, p)

	require.Equal(t, &Handler{
		service: &services.Service{Auth: auth},
		log:     log,
		kafka:   kafka,
		params:  p,
	}, handler)
}
