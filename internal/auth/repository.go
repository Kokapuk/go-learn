package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) createUser(ctx context.Context, params createUserParams) (User, error) {
	var user User

	err := r.db.QueryRow(ctx, `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, created_at
	`, params.Username, params.PasswordHash).Scan(&user.ID, &user.Username, &user.CreatedAt)

	return user, err
}

func (r *Repository) getPublicUserByID(ctx context.Context, userID string) (User, error) {
	var user User

	err := r.db.QueryRow(ctx, `
		SELECT id, username, created_at
		FROM users
		WHERE id = $1
	`, userID).Scan(&user.ID, &user.Username, &user.CreatedAt)

	return user, err
}

func (r *Repository) getUserByUsername(ctx context.Context, username string) (UserWithPassword, error) {
	var user UserWithPassword

	err := r.db.QueryRow(ctx, `
		SELECT id, username, created_at, password_hash
		FROM users
		WHERE username = $1
	`, username).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.PasswordHash)

	return user, err
}
