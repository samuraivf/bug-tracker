package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func KafkaProducer() {
	config := kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "mail",
		BatchTimeout: 1 * time.Millisecond}

	w := kafka.NewWriter(config)

	fmt.Println("Producer configuration: ", config)

	i := 1

	defer func() {
		err := w.Close()
		if err != nil {
			fmt.Println("Error closing producer: ", err)
			return
		}
		fmt.Println("Producer closed")
	}()

	for {
		message := fmt.Sprintf("Message-%d", i)
		err := w.WriteMessages(context.Background(), kafka.Message{Value: []byte(message)})
		if err == nil {
			fmt.Println("Sent message: ", message)
		} else if err == context.Canceled {
			fmt.Println("Context canceled: ", err)
			break
		} else {
			fmt.Println("Error sending message: ", err)
		}
		i++

		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}