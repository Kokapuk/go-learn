package posts

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

func (r *Repository) createPost(ctx context.Context, params createPostParams) (Post, error) {
	var post Post

	err := r.db.QueryRow(ctx, `
		INSERT INTO posts (title, content, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, content, author_id, created_at
	`, params.Title, params.Content, params.AuthorID).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt)

	return post, err
}

func (r *Repository) getPosts(ctx context.Context, params listParams) ([]Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, content, author_id, created_at
		FROM posts
		ORDER BY created_at desc
		LIMIT $1 OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *Repository) getPostsByAuthorID(ctx context.Context, params listParams, authorID string) ([]Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, content, author_id, created_at
		FROM posts
		WHERE author_id = $1
		ORDER BY created_at desc
		LIMIT $2 OFFSET $3
	`, authorID, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
