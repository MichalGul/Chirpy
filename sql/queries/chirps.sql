-- name: CreateChirp :one
INSERT into chirps (
    id, created_at, updated_at, body, user_id
)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps order by created_at asc;

-- name: GetChirpById :one
SELECT * from chirps where chirps.id = $1;

-- name: DeleteChirpById :exec
DELETE from chirps where chirps.id = $1;

-- name: GetChirpsByAuthorId :many
SELECT * from chirps where chirps.user_id = $1 order by created_at asc;