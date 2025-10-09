package service

import (
	"context"
	"errors"
	"fmt"
	"myproject/models"
	"myproject/repository"
	"strings"
)

var (
    ErrInvalidID      = errors.New("invalid note ID")
    ErrNoteNotFound   = errors.New("note not found")
    ErrTitleRequired  = errors.New("title is required")
    ErrTitleTooLong   = errors.New("title too long")
    ErrContentTooLong = errors.New("content too long")
)

type NoteService interface {
	GetNote(ctx context.Context, id int) (*models.Note, error)
	GetAllNotes(ctx context.Context) ([]models.Note, error)
	CreateNote(ctx context.Context, title, content string)(int, error)
	DeleteNote(ctx context.Context, id int) error
}

type noteService struct{
	repo *repository.NoteRepository
}

func CreateNoteService(repo *repository.NoteRepository) *noteService{
	return &noteService{repo: repo}
}

func (s *noteService) GetNote(ctx context.Context,id int) (*models.Note, error){
	
	if id<=0{
		return nil, ErrInvalidID
	}

	note, err := s.repo.GetById(ctx, id)
	if err != nil{
		if errors.Is(err, repository.ErrNotFound) {
            return nil, ErrNoteNotFound
        }
        return nil, fmt.Errorf("service: get-note-by-id id = %d: %w", id, err)
	}

	return note, nil
}

func (s *noteService) GetAllNotes(ctx context.Context) ([]models.Note, error){
	notes, err := s.repo.GetAll(ctx)
	if err != nil{
		return nil, fmt.Errorf("service: get-all-notes: %w", err)
	}

	return notes, nil
}

func (s *noteService) CreateNote(ctx context.Context, title, content string)(int, error){
	if strings.TrimSpace(title) == "" {
		return 0, ErrTitleRequired
	}
	if len(title) > 255 {
		return 0, ErrTitleTooLong
	}
	if len(content) > 5000 {
		return 0, ErrContentTooLong
	}
	
	id, err := s.repo.Create(ctx, title, content)
	if err!=nil{
		return 0, fmt.Errorf("service: create-note: %w", err)
	}

	return id, nil
}

func (s *noteService) DeleteNote(ctx context.Context, id int) error{
	if id <= 0{
		return ErrInvalidID
	}
	
	if err := s.repo.Delete(ctx, id); err != nil{
		if errors.Is(err, repository.ErrNotFound){
			return ErrNoteNotFound
		}
		return fmt.Errorf("service: delete-note id = %d: %w", id, err)
	}
	return nil
}