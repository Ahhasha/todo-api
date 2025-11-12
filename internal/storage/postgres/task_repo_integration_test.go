package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"todo-api/internal/domain"
	"todo-api/internal/storage/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	dsn := "postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable"

	var err error
	db, err = pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	db.Close()
	os.Exit(code)
}

func TestTaskRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := postgres.NewTaskRepo(db)
	_, err := db.Exec(ctx, `
        TRUNCATE TABLE tasks, lists RESTART IDENTITY CASCADE;
    `)
	require.NoError(t, err)

	listID := uuid.New().String()
	_, err = db.Exec(ctx, `
        INSERT INTO lists (id, title, created_at)
        VALUES ($1, 'Integration List', NOW())
        ON CONFLICT (id) DO NOTHING;
    `, listID)
	require.NoError(t, err)

	task := &domain.Task{
		ID:        uuid.New().String(),
		ListID:    listID,
		Text:      "Test integration task",
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repo.Create(ctx, task)
	require.NoError(t, err)
	require.NotEmpty(t, task.ID)

	stored, err := repo.GetByID(ctx, task.ID)
	require.NoError(t, err)
	require.Equal(t, "Test integration task", stored.Text)
	require.False(t, stored.Completed)

	stored.Completed = true
	stored.Text = "Updated text"
	err = repo.Update(ctx, stored)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, task.ID)
	require.NoError(t, err)
	require.True(t, updated.Completed)
	require.Equal(t, "Updated text", updated.Text)

	tasks, total, err := repo.ListByListID(ctx, listID, 10, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, total, 1)
	require.NotEmpty(t, tasks)

	err = repo.Delete(ctx, task.ID)
	require.NoError(t, err)

	deleted, err := repo.GetByID(ctx, task.ID)
	require.Error(t, err)
	require.Nil(t, deleted)
}

func BenchmarkTaskRepository_ListByListID(b *testing.B) {
	ctx := context.Background()
	repo := postgres.NewTaskRepo(db)

	listID := uuid.New().String()
	_, err := db.Exec(ctx, `
		INSERT INTO lists (id, title, created_at)
		VALUES ($1, 'Bench List', NOW())
		ON CONFLICT (id) DO NOTHING;
	`, listID)
	require.NoError(b, err)

	task := &domain.Task{
		ID:        uuid.New().String(),
		ListID:    listID,
		Text:      "Bench task",
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(ctx, task)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := repo.ListByListID(ctx, listID, 100, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}
