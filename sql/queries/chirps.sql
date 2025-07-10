-- name: CreateChirp :one
INSERT INTO chirps (id, created_at,updated_at,body,user_id)
VALUES(
    gen_random_UUID(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: ResetChirpDatabase :exec
DELETE FROM chirps *;