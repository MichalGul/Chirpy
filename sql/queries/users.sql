-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2

)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users where users.email = $1;

-- name: GetUserById :one
SELECT * FROM users where users.id = $1;

-- name: UpdateUser :one
UPDATE users set email = $2, hashed_password = $3, updated_at = NOW() where id = $1 RETURNING *;

-- name: SetChirpyRed :one
UPDATE users set is_chirpy_red = $2, updated_at = NOW() where id = $1 RETURNING *;