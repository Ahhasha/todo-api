package storage

import (
	"context"
	"todo-api/internal/domain"
)

type ListRepository interface {
	Create(ctx context.Context, list *domain.List) (*domain.List, error)
	GetByID(ctx context.Context, id string) (*domain.List, error)
	Update(ctx context.Context, list *domain.List) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*domain.List, int)
	FindWithPagination(ctx context.Context, limit, offset int) ([]*domain.List, int)
	SearchByTitle(ctx context.Context, query string) ([]domain.List, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	ListByListID(ctx context.Context, listID string, limit, offset int) ([]*domain.Task, int, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id string) error
}
