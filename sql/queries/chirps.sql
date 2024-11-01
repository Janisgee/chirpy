-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)

RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps 
ORDER BY created_at ASC;

-- name: GetAllChirpsDesc :many
SELECT * FROM chirps 
ORDER BY created_at DESC;

-- name: GetAllChipsByUserID :many
SELECT * FROM chirps WHERE user_id=$1
ORDER BY created_at ASC;

-- name: GetAllChipsByUserIDDesc :many
SELECT * FROM chirps WHERE user_id=$1
ORDER BY created_at DESC;

-- name: GetOneChirp :one
SELECT * FROM chirps WHERE id=$1 LIMIT 1;


-- name: DeleteOneChirpByUserId :one
DELETE FROM chirps WHERE id=$1 AND user_id=$2
RETURNING *;