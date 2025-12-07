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

func RunUserRegisteredConsumer(ctx context.Context, noteSvc service.NoteService) error {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092"
	}
	topic := os.Getenv("KAFKA_USER_REGISTERED_TOPIC")
	if topic == "" {
		topic = "user_registered"
	}

	fmt.Printf("[KAFKA] starting consumer on broker=%s topic=%s\n", broker, topic)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		Topic:          topic,
		GroupID:        "note-service-consumer",
		CommitInterval: time.Second,
	})

	go func() {
		<-ctx.Done()
		fmt.Println("[KAFKA] context closed, stopping consumer")
		_ = reader.Close()
	}()

	// Основная горутина чтения сообщений Kafka
	go func() {
		for {
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					fmt.Println("[KAFKA] consumer stopped")
					return
				}
				fmt.Printf("[KAFKA] read error: %v\n", err)
				continue
			}

			var ev UserRegisteredEvent
			if err := json.Unmarshal(m.Value, &ev); err != nil {
				fmt.Printf("[KAFKA] failed to unmarshal event: %v\n", err)
				continue
			}

			// ЛОГ №1 — событие получено
			fmt.Printf(
				"[KAFKA] event received topic=%s partition=%d offset=%d user_id=%d email=%s\n",
				m.Topic, m.Partition, m.Offset, ev.UserID, ev.Email,
			)

			// ЛОГ №2 — начинаем создание
			fmt.Printf("[KAFKA] creating welcome note for user=%d\n", ev.UserID)

			title := "Добро пожаловать!"
			content := fmt.Sprintf("Привет, %s! Это ваша первая заметка.", ev.Email)

			id, err := noteSvc.CreateNote(context.Background(), ev.UserID, title, content)
			if err != nil {
				fmt.Printf("[KAFKA] failed to create welcome note for user=%d: %v\n", ev.UserID, err)
				continue
			}

			// ЛОГ №3 — всё успешно
			fmt.Printf("[KAFKA] welcome note created id=%d user=%d\n", id, ev.UserID)
		}
	}()

	return nil
}
