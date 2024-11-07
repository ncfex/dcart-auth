-- name: CreateUser :one
INSERT INTO users (
    id,
    username,
    password_hash,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    $1,
    $2,
    NOW() AT TIME ZONE 'UTC',
    NOW() AT TIME ZONE 'UTC'
) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;