package posts

import (
	"context"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) createPost(ctx context.Context, params createPostParams) (Post, error) {
	return s.repository.createPost(ctx, params)
}

func (s *Service) getPosts(ctx context.Context, params listParams) ([]Post, error) {
	return s.repository.getPosts(ctx, params)
}

func (s *Service) getPostsByAuthorID(ctx context.Context, params listParams, authorID string) ([]Post, error) {
	return s.repository.getPostsByAuthorID(ctx, params, authorID)
}
