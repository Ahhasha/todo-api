package service

import (
	"context"
	"errors"

	"todo-api/internal/domain"
	"todo-api/internal/storage"
)

type ListService interface {
	CreateList(ctx context.Context, title, description string) (*domain.List, error)
	UpdateList(ctx context.Context, id string, title, description string) (*domain.List, error)
	GetAllLists(ctx context.Context) ([]*domain.List, int)
	GetByID(ctx context.Context, id string) (*domain.List, error)
	Delete(ctx context.Context, id string) error
	GetAllListsWithPagination(ctx context.Context, limit, offset int) ([]*domain.List, int)
	SearchByTitle(ctx context.Context, query string) ([]domain.List, error)
}

type listService struct {
	repo storage.ListRepository
}

func NewListService(repo storage.ListRepository) ListService {
	return &listService{
		repo: repo,
	}
}

func (s *listService) CreateList(ctx context.Context, title, description string) (*domain.List, error) {
	if len(title) < 1 || len(title) > 100 {
		return nil, errors.New("title must be 1..100 chars")
	}

	list := domain.NewList(title, description)
	_, _ = s.repo.Create(ctx, list)
	return list, nil
}

func (s *listService) UpdateList(ctx context.Context, id, title, description string) (*domain.List, error) {
	if len(title) < 1 || len(title) > 100 {
		return nil, errors.New("title must be 1..100 chars")
	}

	list, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	list.Title = title
	list.Description = description
	_ = s.repo.Update(ctx, list)
	return list, nil
}

func (s *listService) GetAllLists(ctx context.Context) ([]*domain.List, int) {
	return s.repo.GetAll(ctx)
}

func (s *listService) GetByID(ctx context.Context, id string) (*domain.List, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *listService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *listService) GetAllListsWithPagination(ctx context.Context, limit, offset int) ([]*domain.List, int) {
	return s.repo.FindWithPagination(ctx, limit, offset)
}

func (s *listService) SearchByTitle(ctx context.Context, query string) ([]domain.List, error) {
	return s.repo.SearchByTitle(ctx, query)
}
