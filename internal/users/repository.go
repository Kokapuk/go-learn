package users

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
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
