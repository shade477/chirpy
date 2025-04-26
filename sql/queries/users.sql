-- ':one' indicates the number of rows that it will return
-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email)
VALUES(
    gen_random_uuid(),
    now(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;
