package repository

import (
	"context"
	"todo-app/internal/models"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *models.Todo) error
	GetAll(ctx context.Context) ([]models.Todo, error)
	GetByID(ctx context.Context, id string) (*models.Todo, error)
	Update(ctx context.Context, todo *models.Todo) error
	Delete(ctx context.Context, id string) error
}
