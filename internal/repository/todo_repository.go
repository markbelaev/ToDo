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

func (r *ToDoRepository) GetToDos(ctx context.Context) ([]models.ToDo, error) {
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

func (r *ToDoRepository) GetToDoByID(ctx context.Context, id int64) (*models.ToDo, error) {
	query := "SELECT id, title, description, status, created_at, updated_at FROM todos WHERE id = $1"

	var todo models.ToDo
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		slog.Error("ToDoRepository.GetToDoByID QueryRow:", "error", err)
		return nil, err
	}

	slog.Info("GetToDoByID success")
	return &todo, nil
}

func (r *ToDoRepository) DeleteToDo(ctx context.Context, id int64) error {
	query := "DELETE FROM todos WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		slog.Error("ToDoRepository.DeleteToDo Exec:", "error", err)
		return err
	}
	slog.Info("DeleteToDo success")
	return nil
}

func (r *ToDoRepository) CreateToDo(ctx context.Context, todo *models.ToDo) (*models.ToDo, error) {
	query := "INSERT INTO todos (title, description, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"

	err := r.pool.QueryRow(ctx, query, todo.Title, todo.Description, todo.Status).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		slog.Error("ToDoRepository.CreateToDo Exec:", "error", err)
		return nil, err
	}

	slog.Info("CreateToDo success")
	return todo, nil
}
