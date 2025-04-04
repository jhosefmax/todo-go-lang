// @title Todo API
// @version 1.0
// @description This is a simple TODO list API using Supabase
// @host localhost:8080
// @BasePath /
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/jhosefmoreira/test-go-lang/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	supa "github.com/supabase-community/postgrest-go"
)

// Todo represents a todo item
// @Description Todo item
type Todo struct {
	// The unique identifier for the todo item
	ID string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	// The title of the todo item
	Title string `json:"title" example:"Buy groceries"`
	// A description of what needs to be done
	Description string `json:"description" example:"Get milk and bread"`
	// Whether the todo item has been completed
	Completed bool `json:"completed" example:"false"`
	// When the todo item was created
	CreatedAt time.Time `json:"created_at"`
	// When the todo item was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// @title Todo API
// @version 1.0
// @description This is a simple TODO list API using Supabase
// @host localhost:8080
// @BasePath /
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get Supabase credentials from environment variables
	supabaseUrl := os.Getenv("SUPABASE_URL")
	if supabaseUrl == "" {
		log.Fatal("SUPABASE_URL environment variable is required")
	}

	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")
	if supabaseKey == "" {
		log.Fatal("SUPABASE_ANON_KEY environment variable is required")
	}

	// Initialize Supabase client
	client := supa.NewClient(supabaseUrl+"/rest/v1", "", map[string]string{
		"apikey":        supabaseKey,
		"Authorization": "Bearer " + supabaseKey,
	})

	// Create router
	r := mux.NewRouter()

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // The URL pointing to API definition
	))

	// Routes
	r.HandleFunc("/todos", getTodos(client)).Methods("GET")
	r.HandleFunc("/todos", createTodo(client)).Methods("POST")
	r.HandleFunc("/todos/{id}", updateTodo(client)).Methods("PUT")
	r.HandleFunc("/todos/{id}", deleteTodo(client)).Methods("DELETE")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// getTodos godoc
// @Summary Get all todos
// @Description Get all todo items
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} Todo
// @Router /todos [get]
func getTodos(client *supa.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		resp, count, err := client.From("todos").Select("*", "", false).Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = count // we don't need the count for now

		if err := json.Unmarshal(resp, &todos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(todos)
	}
}

// createTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body Todo true "Todo object"
// @Success 201 {object} Todo
// @Router /todos [post]
func createTodo(client *supa.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Generate UUID for new todo
		todo.ID = uuid.New().String()
		now := time.Now()
		todo.CreatedAt = now
		todo.UpdatedAt = now

		resp, count, err := client.From("todos").Insert(todo, true, "", "", "").Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = count // we don't need the count for now

		var results []Todo
		if err := json.Unmarshal(resp, &results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(results) == 0 {
			http.Error(w, "no todo created", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(results[0])
	}
}

// updateTodo godoc
// @Summary Update a todo
// @Description Update a todo item by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body Todo true "Todo object"
// @Success 200 {object} Todo
// @Router /todos/{id} [put]
func updateTodo(client *supa.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var todo Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		todo.UpdatedAt = time.Now()

		resp, count, err := client.From("todos").Update(todo, "", "").Eq("id", id).Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = count // we don't need the count for now

		var results []Todo
		if err := json.Unmarshal(resp, &results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(results) == 0 {
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(results[0])
	}
}

// deleteTodo godoc
// @Summary Delete a todo
// @Description Delete a todo item by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 204 "No Content"
// @Router /todos/{id} [delete]
func deleteTodo(client *supa.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		resp, count, err := client.From("todos").Delete("", "").Eq("id", id).Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = count // we don't need the count for now
		_ = resp  // we don't need the response for delete

		w.WriteHeader(http.StatusNoContent)
	}
}
