-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_UUID(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;


-- name: ResetUserDatabase :exec
DELETE FROM users *;

-- name: ReturnUserByEmail :one
SELECT * from users
WHERE email = $1;

-- name: UpdateUserData :one
UPDATE users
SET email = $1, hashed_password = $2
where id = $3
RETURNING *;