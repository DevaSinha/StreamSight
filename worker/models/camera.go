package models

import "time"

type Camera struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}
