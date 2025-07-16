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

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpFromID :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: GetChirpsFromAuthor :many
SELECT *
FROM chirps
where user_id = $1
ORDER BY created_at ASC;