package domain

import (
	"time"

	"github.com/google/uuid"
)

type List struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func NewList(title string) *List {
	return &List{
		ID:        uuid.New().String(),
		Title:     title,
		CreatedAt: time.Now(),
	}
}
