package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/adityaladwa/todo-app/internal/db"
    "github.com/go-chi/chi/v5"
    "github.com/sirupsen/logrus"
)

type TodoHandler struct {
    Queries *db.Queries
    Logger  *logrus.Logger
}

type CreateTodoRequest struct {
    Title string `json:"title"`
}

type UpdateTodoRequest struct {
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    var req CreateTodoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.Logger.Error("Invalid request payload", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    todo, err := h.Queries.CreateTodo(r.Context(), db.CreateTodoParams{
        Title:     req.Title,
        Completed: false,
    })
    if err != nil {
        h.Logger.Error("Failed to create todo", err)
        http.Error(w, "Failed to create todo", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        h.Logger.Error("Invalid todo ID", err)
        http.Error(w, "Invalid todo ID", http.StatusBadRequest)
        return
    }

    todo, err := h.Queries.GetTodo(r.Context(), int32(id))
    if err != nil {
        h.Logger.Error("Todo not found", err)
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
    limitStr := r.URL.Query().Get("limit")
    offsetStr := r.URL.Query().Get("offset")

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }

    offset, err := strconv.Atoi(offsetStr)
    if err != nil || offset < 0 {
        offset = 0
    }

    todos, err := h.Queries.ListTodos(
        r.Context(),
        db.ListTodosParams{
            int32(limit),
            int32(offset),
        },
    )
    if err != nil {
        h.Logger.Error("Failed to list todos", err)
        http.Error(w, "Failed to list todos", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        h.Logger.Error("Invalid todo ID", err)
        http.Error(w, "Invalid todo ID", http.StatusBadRequest)
        return
    }

    var req UpdateTodoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.Logger.Error("Invalid request payload", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    todo, err := h.Queries.UpdateTodo(r.Context(), db.UpdateTodoParams{
        Title:     req.Title,
        Completed: req.Completed,
        ID:        int32(id),
    })
    if err != nil {
        h.Logger.Error("Failed to update todo", err)
        http.Error(w, "Failed to update todo", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        h.Logger.Error("Invalid todo ID", err)
        http.Error(w, "Invalid todo ID", http.StatusBadRequest)
        return
    }

    err = h.Queries.DeleteTodo(r.Context(), int32(id))
    if err != nil {
        h.Logger.Error("Failed to delete todo", err)
        http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
