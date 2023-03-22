package kafka

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
)

type KafkaWriter struct {
	writer *kafkago.Writer
	log    log.Log
}

type Kafka interface {
	Close() error
	Write(message string) error
}

func NewKafkaWriter(config kafkago.WriterConfig, log log.Log) *KafkaWriter {
	return &KafkaWriter{
		writer: kafkago.NewWriter(config),
		log: log,
	}
}

func (w *KafkaWriter) Close() error {
	return w.writer.Close()
}

func (w *KafkaWriter) Write(message string) error {
	return w.writer.WriteMessages(
		context.Background(),
		kafkago.Message{
			Value: []byte(message),
		},
	)
}
