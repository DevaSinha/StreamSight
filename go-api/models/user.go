package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hide password in JSON
	CreatedAt time.Time `json:"created_at"`
}
