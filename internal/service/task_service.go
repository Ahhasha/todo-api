package service

import (
	"context"
	"errors"
	"time"

	"todo-api/internal/domain"
	"todo-api/internal/storage"

	"github.com/google/uuid"
)

type TaskService struct {
	repo     storage.TaskRepository
	listRepo storage.ListRepository
}

func NewTaskService(repo storage.TaskRepository, listRepo storage.ListRepository) *TaskService {
	return &TaskService{repo: repo, listRepo: listRepo}
}

func (s *TaskService) CreateTask(ctx context.Context, listID, text string) (*domain.Task, error) {
	if len(text) < 1 || len(text) > 500 {
		return nil, errors.New("text must be 1..500 chars")
	}

	if _, err := s.listRepo.GetByID(ctx, listID); err != nil {
		return nil, err
	}

	task := &domain.Task{
		ID:        uuid.NewString(),
		ListID:    listID,
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context, listID string, limit, offset int) ([]*domain.Task, int, error) {
	if _, err := s.listRepo.GetByID(ctx, listID); err != nil {
		return nil, 0, err
	}
	return s.repo.ListByListID(ctx, listID, limit, offset)
}

func (s *TaskService) UpdateTask(ctx context.Context, id, text string, completed *bool) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if text != "" {
		if len(text) < 1 || len(text) > 500 {
			return nil, errors.New("text must be 1..500 chars")
		}
		task.Text = text
	}

	if completed != nil {
		task.Completed = *completed
	}

	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
