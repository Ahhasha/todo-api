package mem

import (
	"context"
	"errors"
	"sync"
	"todo-api/internal/domain"
)

type ListRepo interface {
	Create(ctx context.Context, list *domain.List) error
	GetAll(ctx context.Context) ([]*domain.List, int)
	GetByID(ctx context.Context, id string) (*domain.List, error)
	Update(ctx context.Context, list *domain.List) error
	Delete(ctx context.Context, id string) error
	FindWithPagination(ctx context.Context, limit, offset int) ([]*domain.List, int)
}

type listRepo struct {
	mu    sync.RWMutex
	lists map[string]*domain.List
}

func NewListRepo() ListRepo {
	return &listRepo{
		lists: make(map[string]*domain.List),
	}
}

func (r *listRepo) Create(ctx context.Context, list *domain.List) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.lists[list.ID] = list
		return nil
	}
}
func (r *listRepo) GetAll(ctx context.Context) ([]*domain.List, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, 0
	default:
		var result []*domain.List
		for _, list := range r.lists {
			result = append(result, list)
		}
		return result, len(result)
	}
}
func (r *listRepo) GetByID(ctx context.Context, id string) (*domain.List, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		list, ok := r.lists[id]
		if !ok {
			return nil, errors.New("list not found")
		}
		return list, nil
	}
}
func (r *listRepo) Update(ctx context.Context, list *domain.List) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, ok := r.lists[list.ID]
		if !ok {
			return errors.New("list not found")
		}
		r.lists[list.ID] = list
		return nil
	}
}
func (r *listRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, ok := r.lists[id]
		if !ok {
			return errors.New("list not found")
		}
		delete(r.lists, id)
		return nil
	}
}

func (r *listRepo) FindWithPagination(ctx context.Context, limit, offset int) ([]*domain.List, int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, 0
	default:
		var result []*domain.List
		for _, list := range r.lists {
			result = append(result, list)
		}
		total := len(result)

		if offset > total {
			return []*domain.List{}, total
		}
		result = result[offset:]

		if limit > 0 && limit < len(result) {
			result = result[:limit]
		}
		return result, total
	}

}
