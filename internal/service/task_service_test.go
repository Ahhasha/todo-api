package service_test

import (
	"context"
	"errors"
	"testing"

	"todo-api/internal/domain"
	"todo-api/internal/service"
)

type mockTaskRepo struct {
	createFunc func(ctx context.Context, task *domain.Task) error
}

func (m *mockTaskRepo) Create(ctx context.Context, task *domain.Task) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, task)
	}
	return nil
}

func (m *mockTaskRepo) GetByID(ctx context.Context, id string) (*domain.Task, error) { return nil, nil }

func (m *mockTaskRepo) ListByListID(ctx context.Context, listID string, limit, offset int) ([]*domain.Task, int, error) {
	return nil, 0, nil
}

func (m *mockTaskRepo) Update(ctx context.Context, task *domain.Task) error { return nil }

func (m *mockTaskRepo) Delete(ctx context.Context, id string) error { return nil }

type mockListRepo struct {
	getByIDFunc func(ctx context.Context, id string) (*domain.List, error)
}

func (m *mockListRepo) Create(ctx context.Context, list *domain.List) (*domain.List, error) {
	return list, nil
}

func (m *mockListRepo) GetByID(ctx context.Context, id string) (*domain.List, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return &domain.List{ID: id, Title: "Test List"}, nil
}

func (m *mockListRepo) Update(ctx context.Context, list *domain.List) error { return nil }

func (m *mockListRepo) Delete(ctx context.Context, id string) error { return nil }

func (m *mockListRepo) GetAll(ctx context.Context) ([]*domain.List, int) {
	return []*domain.List{}, 0
}

func (m *mockListRepo) FindWithPagination(ctx context.Context, limit, offset int) ([]*domain.List, int) {
	return []*domain.List{}, 0
}

func (m *mockListRepo) CreateWithItems(ctx context.Context, title string, items []string) error {
	return nil
}

func (m *mockListRepo) SearchByTitle(ctx context.Context, query string) ([]domain.List, error) {
	return []domain.List{}, nil
}

func TestTaskService_CreateTask_Success(t *testing.T) {
	taskRepo := &mockTaskRepo{
		createFunc: func(ctx context.Context, task *domain.Task) error {
			task.ID = "task-1"
			return nil
		},
	}

	listRepo := &mockListRepo{
		getByIDFunc: func(ctx context.Context, id string) (*domain.List, error) {
			return &domain.List{ID: id, Title: "My list"}, nil
		},
	}

	svc := service.NewTaskService(taskRepo, listRepo)

	task, err := svc.CreateTask(context.Background(), "list-1", "Купить хлеб")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.Text != "Купить хлеб" {
		t.Errorf("expected text 'Купить хлеб', got '%s'", task.Text)
	}
}

func TestTaskService_CreateTask_ValidationError(t *testing.T) {
	svc := service.NewTaskService(&mockTaskRepo{}, &mockListRepo{})

	_, err := svc.CreateTask(context.Background(), "list-1", "")
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
}

func TestTaskService_CreateTask_ListNotFound(t *testing.T) {
	taskRepo := &mockTaskRepo{}
	listRepo := &mockListRepo{
		getByIDFunc: func(ctx context.Context, id string) (*domain.List, error) {
			return nil, errors.New("list not found")
		},
	}

	svc := service.NewTaskService(taskRepo, listRepo)

	_, err := svc.CreateTask(context.Background(), "missing-list", "Задача")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
