package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

 	"golang.org/x/crypto/bcrypt"

	"user-service/models"
	"user-service/repository"
)

var (
    ErrEmailRequired     = errors.New("email is required")
    ErrEmailInvalid      = errors.New("email is invalid")
    ErrPasswordRequired  = errors.New("password is required")
    ErrPasswordTooShort  = errors.New("password is too short")
    ErrEmailAlreadyTaken = errors.New("email already registered")

    ErrInvalidCredentials = errors.New("invalid email or password")
)


type UserService interface {
    RegisterUser(ctx context.Context, email, password string) (models.User, error)
    LoginUser(ctx context.Context, email, password string) (models.User, error)
}


type userService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, email, password string) (models.User, error) {
    email = strings.TrimSpace(email)
    if email == "" {
        return models.User{}, ErrEmailRequired
    }
    if !strings.Contains(email, "@") {
        return models.User{}, ErrEmailInvalid
    }
    if len(password) == 0 {
        return models.User{}, ErrPasswordRequired
    }
    if len(password) < 6 {
        return models.User{}, ErrPasswordTooShort
    }

    // проверим, нет ли уже такого email
    _, err := s.repo.GetByEmail(ctx, email)
    if err == nil {
        // нашли пользователя
        return models.User{}, ErrEmailAlreadyTaken
    }
    if !errors.Is(err, repository.ErrNotFound) {
        // какая-то другая ошибка из репозитория
        return models.User{}, fmt.Errorf("service: check-email: %w", err)
    }

    // хэшируем пароль
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return models.User{}, fmt.Errorf("service: hash-password: %w", err)
    }

    // создаём пользователя с хэшированным паролем
    id, err := s.repo.Create(ctx, email, string(hash))
    if err != nil {
        return models.User{}, fmt.Errorf("service: create-user: %w", err)
    }

    return models.User{
        Id:       id,
        Email:    email,
        Password: "", // наружу пароль не отдаём
    }, nil
}

func (s *userService) LoginUser(ctx context.Context, email, password string) (models.User, error) {
    email = strings.TrimSpace(email)
    if email == "" || password == "" {
        return models.User{}, ErrInvalidCredentials
    }

    u, err := s.repo.GetByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            return models.User{}, ErrInvalidCredentials
        }
        return models.User{}, fmt.Errorf("service: login-get-by-email: %w", err)
    }

    // сравниваем хэш и введённый пароль
    if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
        return models.User{}, ErrInvalidCredentials
    }

    return u, nil
}

