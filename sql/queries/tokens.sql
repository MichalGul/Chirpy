-- name: CreateRefreshToken :one
INSERT into refresh_tokens (
    token, created_at, updated_at, user_id, expires_at, revoked_at
)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * from refresh_tokens where refresh_tokens.token = $1;

-- name: GetUserFromRefreshToken :one
SELECT user_id from refresh_tokens where refresh_tokens.token = $1;

-- name: SetRevokeOnToken :one
UPDATE refresh_tokens set revoked_at = NOW(), updated_at = NOW() where token=$1 RETURNING *;