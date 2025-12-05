package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"myproject/models"
)

type NotesCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewNotesCache() *NotesCache {
	addr := getEnv("REDIS_ADDR", "redis:6379")
	dbStr := getEnv("REDIS_DB", "0")

	dbNum, err := strconv.Atoi(dbStr)
	if err != nil {
		dbNum = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   dbNum,
	})

	return &NotesCache{
		client: client,
		ttl:    90 * time.Second, // кэш живёт 90 сек
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (c *NotesCache) key(userID int) string {
	return fmt.Sprintf("notes:%d", userID)
}

// Пытаемся взять заметки из кэша.
// Возвращаем: данные, найдено ли в кэше, ошибка.
func (c *NotesCache) GetNotes(ctx context.Context, userID int) ([]models.Note, bool, error) {
	if c == nil || c.client == nil {
		return nil, false, nil
	}

	data, err := c.client.Get(ctx, c.key(userID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil // в кэше нет
		}
		return nil, false, fmt.Errorf("redis get: %w", err)
	}

	var notes []models.Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return nil, false, fmt.Errorf("unmarshal cached notes: %w", err)
	}

	return notes, true, nil
}

// Кладём заметки в кэш
func (c *NotesCache) SetNotes(ctx context.Context, userID int, notes []models.Note) error {
	if c == nil || c.client == nil {
		return nil
	}

	data, err := json.Marshal(notes)
	if err != nil {
		return fmt.Errorf("marshal notes: %w", err)
	}

	if err := c.client.Set(ctx, c.key(userID), data, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

// Инвалидируем (удаляем) кэш для пользователя — после create/update/delete
func (c *NotesCache) Invalidate(ctx context.Context, userID int) error {
	if c == nil || c.client == nil {
		return nil
	}
	if err := c.client.Del(ctx, c.key(userID)).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}
