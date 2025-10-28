package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"todo-api/internal/domain"
)

var ErrNotFound = errors.New("not found")

type ListRepo struct {
	pool *pgxpool.Pool
}

func NewListRepo(pool *pgxpool.Pool) *ListRepo {
	return &ListRepo{pool: pool}
}

func (r *ListRepo) Create(ctx context.Context, list *domain.List) (*domain.List, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO lists (id, title, description)
        VALUES ($1, $2, $3)
        RETURNING created_at
    `
	err := r.pool.QueryRow(ctx, query, list.ID, list.Title, list.Description).Scan(&list.CreatedAt)
	if err != nil {
		return list, fmt.Errorf("create list: %w", err)
	}

	return list, nil
}

func (r *ListRepo) GetByID(ctx context.Context, id string) (*domain.List, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
        SELECT id, title, description, created_at
        FROM lists
        WHERE id = $1
    `
	var list domain.List
	err := r.pool.QueryRow(ctx, query, id).Scan(&list.ID, &list.Title, &list.Description, &list.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get list by id: %w", err)
	}
	return &list, nil
}

func (r *ListRepo) Update(ctx context.Context, list *domain.List) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
        UPDATE lists
        SET title = $2, description = $3
        WHERE id = $1
        RETURNING id, title, description, created_at
    `
	return r.pool.QueryRow(ctx, query, list.ID, list.Title, list.Description).Scan(&list.ID, &list.Title, &list.Description, &list.CreatedAt)
}

func (r *ListRepo) Delete(ctx context.Context, id string) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM lists WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete list: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ListRepo) GetAll(ctx context.Context) ([]*domain.List, int) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `SELECT id, title, description, created_at FROM lists ORDER BY created_at DESC`)
	if err != nil {
		return nil, 0
	}
	defer rows.Close()

	var lists []*domain.List
	for rows.Next() {
		var list domain.List
		if err := rows.Scan(&list.ID, &list.Title, &list.Description, &list.CreatedAt); err == nil {
			lists = append(lists, &list)
		}
	}
	return lists, len(lists)
}

func (r *ListRepo) FindWithPagination(ctx context.Context, limit, offset int) ([]*domain.List, int) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx, `
        SELECT id, title, description, created_at
        FROM lists
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `, limit, offset)
	if err != nil {
		return nil, 0
	}
	defer rows.Close()

	var lists []*domain.List
	for rows.Next() {
		var list domain.List
		if err := rows.Scan(&list.ID, &list.Title, &list.Description, &list.CreatedAt); err == nil {
			lists = append(lists, &list)
		}
	}

	var total int
	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM lists`).Scan(&total)
	return lists, total
}

func (r *ListRepo) CreateWithItems(ctx context.Context, title string, items []string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	listID := uuid.New()
	_, err = tx.Exec(ctx, "INSERT INTO lists (id, title) VALUES ($1, $2)", listID, title)
	if err != nil {
		return fmt.Errorf("create list: %w", err)
	}

	for _, item := range items {
		_, err = tx.Exec(ctx, "INSERT INTO items (list_id, text) VALUES ($1, $2)", listID, item)
		if err != nil {
			return fmt.Errorf("create item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (r *ListRepo) SearchByTitle(ctx context.Context, query string) ([]domain.List, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sqlQuery := `
		SELECT id, title, description, created_at
		FROM lists
		WHERE title ILIKE $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("search lists by title: %w", err)
	}
	defer rows.Close()

	var lists []domain.List
	for rows.Next() {
		var list domain.List
		if err := rows.Scan(&list.ID, &list.Title, &list.Description, &list.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan list: %w", err)
		}
		lists = append(lists, list)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return lists, nil
}
