# Todo App with Go and Supabase

A simple TODO list application built with Go and Supabase for data storage. The application follows clean architecture principles, making it easy to switch between different database providers.

## Prerequisites

1. Install Go (1.16 or later)
   - Visit https://golang.org/dl/
   - Download and install for your operating system

2. Supabase Setup
   - Create a Supabase project at https://supabase.com
   - Get your project URL and API Key
   - Copy `.env.example` to `.env` and update with your credentials

## Project Structure

```
todo-app/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handlers/
│   │   └── todo.go            # HTTP handlers
│   ├── models/
│   │   └── todo.go            # Data models
│   ├── repository/
│   │   └── todo_repository.go # Repository interface
│   └── database/
│       └── supabase/
│           └── todo_repository.go # Supabase implementation
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Database Setup

1. In your Supabase project, create a new table called `todos` with the following structure:

```sql
create table todos (
  id uuid default uuid_generate_v4() primary key,
  title text not null,
  description text,
  completed boolean default false,
  created_at timestamp with time zone default timezone('utc'::text, now()),
  updated_at timestamp with time zone default timezone('utc'::text, now())
);
```

## Configuration

1. Copy the example environment file:
   ```
   cp .env.example .env
   ```

2. Update `.env` with your Supabase credentials:
   ```
   DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.xxxx.supabase.co:5432/postgres
   PORT=8080
   ```

## Running the Application

1. Install dependencies:
   ```
   go mod tidy
   ```

2. Run the application:
   ```
   go run cmd/main.go
   ```

The server will start at `http://localhost:8080`

## API Endpoints

- `GET /todos` - List all todos
- `POST /todos` - Create a new todo
- `PUT /todos/:id` - Update a todo
- `DELETE /todos/:id` - Delete a todo

## Example Requests

### Create a Todo
```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Get milk and bread"}'
```

### List All Todos
```bash
curl http://localhost:8080/todos
```

### Update a Todo
```bash
curl -X PUT http://localhost:8080/todos/{id} \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Get milk and bread","completed":true}'
```

### Delete a Todo
```bash
curl -X DELETE http://localhost:8080/todos/{id}
```

## Adding New Database Providers

To add support for a different database, create a new implementation of the `TodoRepository` interface in the `internal/database` directory. The interface is defined in `internal/repository/todo_repository.go`.
