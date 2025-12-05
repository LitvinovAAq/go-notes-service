package events

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type UserRegisteredEvent struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	// можно добавить created_at, но для задачи не обязательно
}

func NewUserRegisteredWriter() *kafka.Writer {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}
	topic := os.Getenv("KAFKA_USER_REGISTERED_TOPIC")
	if topic == "" {
		topic = "user_registered"
	}

	return &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
}

func PublishUserRegistered(ctx context.Context, w *kafka.Writer, userID int, email string) error {
	ev := UserRegisteredEvent{
		UserID: userID,
		Email:  email,
	}
	data, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprint(userID)),
		Value: data,
		Time:  time.Now(),
	}

	if err := w.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}

	return nil
}
