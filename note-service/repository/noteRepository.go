package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"myproject/models"
)

var ErrNotFound = errors.New("not found")

type NoteRepo interface {
	GetAll(ctx context.Context, userID int) ([]models.Note, error)
	GetById(ctx context.Context, userID, id int) (models.Note, error)
	Create(ctx context.Context, userID int, title, content string) (int, error)
	Delete(ctx context.Context, userID, id int) error
	Update(ctx context.Context, userID, id int, title *string, content *string) (models.Note, error)
}

type NoteRepository struct {
	db *sql.DB
}

func CreateNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// Получить все заметки конкретного пользователя
func (r *NoteRepository) GetAll(ctx context.Context, userID int) ([]models.Note, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, title, content FROM notes WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("repo: get all notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.Id, &note.UserID, &note.Title, &note.Content)
		if err != nil {
			return nil, fmt.Errorf("repo: scan notes %w", err)
		}
		notes = append(notes, note)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: rows: %w", err)
	}
	return notes, nil
}

// Получить одну заметку пользователя по id
func (r *NoteRepository) GetById(ctx context.Context, userID, id int) (models.Note, error) {
	var note models.Note
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, title, content FROM notes WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&note.Id, &note.UserID, &note.Title, &note.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Note{}, ErrNotFound
		}
		return models.Note{}, fmt.Errorf("repo: get note by id: %w", err)
	}
	return note, nil
}

// Создать заметку для пользователя
func (r *NoteRepository) Create(ctx context.Context, userID int, title, content string) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id`,
		userID, title, content,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo: create note title = %q: %w", title, err)
	}

	return id, nil
}

// Удалить заметку пользователя
func (r *NoteRepository) Delete(ctx context.Context, userID, id int) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM notes WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("repo: delete note id=%d: %w", id, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo: delete id=%d: rowsAffected: %w", id, err)
	}
	if n == 0 {
		return ErrNotFound
	}

	return nil
}

// Частично обновить заметку пользователя
func (r *NoteRepository) Update(ctx context.Context, userID, id int, title *string, content *string) (models.Note, error) {
	// Проверим, есть ли такая заметка и принадлежит ли она этому пользователю
	existing, err := r.GetById(ctx, userID, id)
	if err != nil {
		return models.Note{}, err // здесь ErrNotFound пробросится вверх
	}

	// Обновляем только те поля, которые пришли
	newTitle := existing.Title
	newContent := existing.Content

	if title != nil {
		newTitle = *title
	}
	if content != nil {
		newContent = *content
	}

	query := `UPDATE notes SET title = $1, content = $2 WHERE id = $3 AND user_id = $4`
	_, err = r.db.ExecContext(ctx, query, newTitle, newContent, id, userID)
	if err != nil {
		return models.Note{}, fmt.Errorf("repo: update-note: %w", err)
	}

	// Возвращаем обновлённую заметку
	return models.Note{
		Id:      id,
		UserID:  userID,
		Title:   newTitle,
		Content: newContent,
	}, nil
}
