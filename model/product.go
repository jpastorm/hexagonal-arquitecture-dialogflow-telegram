package model

import (
	"time"
)

// User model
type Product struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Price     float32   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HasID returns true if the user has a valid id
func (p *Product) HasID() bool { return p.ID > 0 }

// users slice of user
type Products []Product

// IsEmpty returns true if users is empty
func (p Products) IsEmpty() bool { return len(p) == 0 }
