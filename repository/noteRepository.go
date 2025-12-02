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
	GetAll(ctx context.Context) ([]models.Note, error)
	GetByID(ctx context.Context, id int) (*models.Note, error)
	Create(ctx context.Context, title, content string) (int, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, title *string, content *string) (models.Note, error)
}

type NoteRepository struct {
	db *sql.DB
}

func CreateNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) GetAll(ctx context.Context) ([]models.Note, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, content FROM notes")
	if err != nil {
		return nil, fmt.Errorf("repo: get all notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.Id, &note.Title, &note.Content)
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

func (r *NoteRepository) GetById(ctx context.Context, id int) (*models.Note, error) {
	var note models.Note
	err := r.db.QueryRowContext(ctx, "SELECT id, title, content FROM notes WHERE id = $1",
		id).Scan(&note.Id, &note.Title, &note.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get note by id: %w", err)
	}
	return &note, nil
}

func (r *NoteRepository) Create(ctx context.Context, title, content string) (int, error) {

	var note models.Note
	err := r.db.QueryRowContext(ctx, "INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id", title, content).
		Scan(&note.Id)
	if err != nil {
		return 0, fmt.Errorf("repo: create note title = %q: %w", title, err)
	}

	return note.Id, nil
}

func (r *NoteRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM notes WHERE id = $1", id)
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

func (r *NoteRepository) Update(ctx context.Context, id int, title *string, content *string) (models.Note, error) {
	// Проверим, есть ли такая заметка
	existing, err := r.GetById(ctx, id)
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

	query := `UPDATE notes SET title = $1, content = $2 WHERE id = $3`
	_, err = r.db.ExecContext(ctx, query, newTitle, newContent, id)
	if err != nil {
		return models.Note{}, fmt.Errorf("repo: update-note: %w", err)
	}

	// Возвращаем обновлённую заметку
	return models.Note{
		Id:      id,
		Title:   newTitle,
		Content: newContent,
	}, nil
}
