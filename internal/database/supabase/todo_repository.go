package supabase

import (
	"context"
	"fmt"
	"time"
	"todo-app/internal/models"

	"github.com/supabase-community/supabase-go"
)

type TodoRepository struct {
	client *supabase.Client
}

func NewTodoRepository(supabaseURL, supabaseKey string) (*TodoRepository, error) {
	client := supabase.CreateClient(supabaseURL, supabaseKey)
	return &TodoRepository{
		client: client,
	}, nil
}

func (r *TodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	var result models.Todo
	err := r.client.DB.From("todos").Insert(todo).Single().Execute(&result)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	
	todo.ID = result.ID
	return nil
}

func (r *TodoRepository) GetAll(ctx context.Context) ([]models.Todo, error) {
	var todos []models.Todo
	err := r.client.DB.From("todos").Select("*").Execute(&todos)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}
	return todos, nil
}

func (r *TodoRepository) GetByID(ctx context.Context, id string) (*models.Todo, error) {
	var todo models.Todo
	err := r.client.DB.From("todos").Select("*").Eq("id", id).Single().Execute(&todo)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	return &todo, nil
}

func (r *TodoRepository) Update(ctx context.Context, todo *models.Todo) error {
	todo.UpdatedAt = time.Now()
	err := r.client.DB.From("todos").Update(todo).Eq("id", todo.ID).Execute(nil)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	err := r.client.DB.From("todos").Delete().Eq("id", id).Execute(nil)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}
