package models

import "time"

type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // не отдаём наружу
	CreatedAt time.Time `json:"created_at"`
}
