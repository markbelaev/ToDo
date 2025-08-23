package repository

import (
	"GIN/internal/models"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ToDoRepository struct {
	pool *pgxpool.Pool
}

func NewToDoRepository(pool *pgxpool.Pool) *ToDoRepository {
	return &ToDoRepository{
		pool: pool,
	}
}

func (r *ToDoRepository) GetAllToDos(ctx context.Context) ([]models.ToDo, error) {
	query := "SELECT id, title, description, status, created_at, updated_at FROM todos"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		slog.Error("ToDoRepository.GetAllToDos Query:", "error", err)
		return nil, err
	}
	defer rows.Close()

	var todos []models.ToDo
	for rows.Next() {
		var todo models.ToDo
		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Status,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			slog.Error("ToDoRepository.GetAllToDos Scan:", "error", err)
			return nil, err
		}

		todos = append(todos, todo)
	}

	slog.Info("GetAllToDos success")
	return todos, nil
}
