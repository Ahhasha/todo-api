package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"todo-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *taskRepo {
	return &taskRepo{pool: pool}
}

func (r *taskRepo) Create(ctx context.Context, t *domain.Task) error {
	query := `INSERT INTO tasks (id, list_id, text, completed)
	          VALUES ($1,$2,$3,$4)
	          RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, query, t.ID, t.ListID, t.Text, t.Completed).
		Scan(&t.CreatedAt, &t.UpdatedAt)
}

func (r *taskRepo) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	var t domain.Task
	err := r.pool.QueryRow(ctx, `SELECT id, list_id, text, completed, created_at, updated_at FROM tasks WHERE id=$1`, id).
		Scan(&t.ID, &t.ListID, &t.Text, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &t, err
}

func (r *taskRepo) ListByListID(ctx context.Context, listID string, limit, offset int) ([]*domain.Task, int, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, list_id, text, completed, created_at, updated_at
		 FROM tasks WHERE list_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, listID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var res []*domain.Task
	for rows.Next() {
		var t domain.Task
		_ = rows.Scan(&t.ID, &t.ListID, &t.Text, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
		res = append(res, &t)
	}
	var total int
	_ = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tasks WHERE list_id=$1`, listID).Scan(&total)
	return res, total, rows.Err()
}

func (r *taskRepo) Update(ctx context.Context, t *domain.Task) error {
	t.UpdatedAt = time.Now().UTC()
	query := `UPDATE tasks SET text=$2, completed=$3, updated_at=$4 WHERE id=$1`
	ct, err := r.pool.Exec(ctx, query, t.ID, t.Text, t.Completed, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *taskRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id=$1`, id)
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}
