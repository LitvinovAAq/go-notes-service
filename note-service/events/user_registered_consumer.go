package events

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"

	"myproject/service"
)

type UserRegisteredEvent struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

// Запускаем отдельную горутину, которая слушает Kafka и создаёт приветственные заметки
func RunUserRegisteredConsumer(ctx context.Context, noteSvc service.NoteService) error {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}
	topic := os.Getenv("KAFKA_USER_REGISTERED_TOPIC")
	if topic == "" {
		topic = "user_registered"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		Topic:          topic,
		GroupID:        "note-service-consumer",
		CommitInterval: time.Second,
	})
	// reader.Close() делаем через горутину/контекст при завершении
	go func() {
		<-ctx.Done()
		_ = reader.Close()
	}()

	go func() {
		for {
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				// при отменённом контексте выходим
				if ctx.Err() != nil {
					return
				}
				fmt.Printf("kafka read error: %v\n", err)
				continue
			}

			var ev UserRegisteredEvent
			if err := json.Unmarshal(m.Value, &ev); err != nil {
				fmt.Printf("failed to unmarshal user_registered: %v\n", err)
				continue
			}

			// создаём приветственную заметку
			title := "Добро пожаловать!"
			content := fmt.Sprintf("Привет, %s! Это ваша первая заметка.", ev.Email)

			_, err = noteSvc.CreateNote(context.Background(), ev.UserID, title, content)
			if err != nil {
				fmt.Printf("failed to create welcome note for user %d: %v\n", ev.UserID, err)
				continue
			}

			fmt.Printf("welcome note created for user %d\n", ev.UserID)
		}
	}()

	return nil
}
