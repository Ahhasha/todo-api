package domain

import (
	"time"

	"github.com/google/uuid"
)

type List struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewList(title, description string) *List {
	return &List{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
