// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_UUID(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const resetUserDatabase = `-- name: ResetUserDatabase :exec
DELETE FROM users *
`

func (q *Queries) ResetUserDatabase(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetUserDatabase)
	return err
}

const returnUserByEmail = `-- name: ReturnUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red from users
WHERE email = $1
`

func (q *Queries) ReturnUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, returnUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUserData = `-- name: UpdateUserData :one
UPDATE users
SET email = $1, hashed_password = $2
where id = $3
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type UpdateUserDataParams struct {
	Email          string
	HashedPassword string
	ID             uuid.UUID
}

func (q *Queries) UpdateUserData(ctx context.Context, arg UpdateUserDataParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserData, arg.Email, arg.HashedPassword, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const upgrateToChirpyRed = `-- name: UpgrateToChirpyRed :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
`

func (q *Queries) UpgrateToChirpyRed(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, upgrateToChirpyRed, id)
	return err
}
