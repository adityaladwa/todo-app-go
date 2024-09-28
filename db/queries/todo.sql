-- name: CreateTodo :one
INSERT INTO todos (title, completed)
VALUES ($1, $2)
RETURNING id, title, completed, created_at;

-- name: GetTodo :one
SELECT id, title, completed, created_at
FROM todos
WHERE id = $1;

-- name: ListTodos :many
SELECT id, title, completed, created_at
FROM todos
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTodo :one
UPDATE todos
SET title = $1, completed = $2
WHERE id = $3
RETURNING id, title, completed, created_at;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;
