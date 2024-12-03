-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW() AT TIME ZONE 'UTC',
    NOW() AT TIME ZONE 'UTC',
    $2,
    $3
)
RETURNING *;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens SET
    revoked_at = NOW() AT TIME ZONE 'UTC',
    updated_at = NOW() AT TIME ZONE 'UTC'
WHERE token = $1
RETURNING *;

-- name: GetTokenByTokenString :one
SELECT * FROM refresh_tokens
WHERE token = $1
AND revoked_at IS NULL
AND expires_at > NOW() AT TIME ZONE 'UTC';

-- name: SaveToken :exec
UPDATE refresh_tokens
SET
    user_id = $2,
    created_at = $3,
    updated_at = $4,
    expires_at = $5,
    revoked_at = $6
WHERE token = $1;
