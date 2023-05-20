package handler

import (
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/kafka"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
)

type Handler struct {
	service *services.Service
	log     log.Log
	kafka   kafka.Kafka
	params  Params
}

func NewHandler(s *services.Service, log log.Log, kafkaWriter kafka.Kafka, params Params) *Handler {
	return &Handler{s, log, kafkaWriter, params}
}
