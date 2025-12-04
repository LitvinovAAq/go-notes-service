package repository

import (
	"context"
	"database/sql"
	"fmt"

	"user-service/models"
)

var ErrNotFound = fmt.Errorf("user not found")

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, email, password string) (int, error) {
	var id int
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	err := r.DB.QueryRowContext(ctx, query, email, password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo: create-user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&u.Id, &u.Email, &u.Password, &u.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNotFound
		}
		return models.User{}, fmt.Errorf("repo: get-by-email: %w", err)
	}
	return u, nil
}
