package model

import (
	"time"
)

// User model
type User struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	TelegramID string    `json:"telegram_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// HasID returns true if the user has a valid id
func (c *User) HasID() bool { return c.ID > 0 }

// users slice of user
type Users []User

// IsEmpty returns true if users is empty
func (c Users) IsEmpty() bool { return len(c) == 0 }
