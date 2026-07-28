package posts

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repository *Repository
	cache      *Cache
}

func NewService(repository *Repository, cache *Cache) *Service {
	return &Service{repository: repository, cache: cache}
}

func (s *Service) createPost(ctx context.Context, params createPostParams) (Post, error) {
	post, err := s.repository.createPost(ctx, params)
	if err != nil {
		return Post{}, err
	}

	if err := s.cache.invalidatePosts(ctx); err != nil {
		log.Println(err)
	}

	return post, nil
}

func (s *Service) getPosts(ctx context.Context, params listParams) ([]Post, error) {
	cachedPosts, err := s.cache.getPosts(ctx, params)
	if err == nil {
		return cachedPosts, nil
	}
	if !errors.Is(err, redis.Nil) {
		log.Println(err)
	}

	posts, err := s.repository.getPosts(ctx, params)
	if err != nil {
		return nil, err
	}

	err = s.cache.cachePosts(ctx, params, posts)
	if err != nil {
		log.Println(err)
	}

	return posts, nil
}

func (s *Service) getPostsByAuthorID(ctx context.Context, params listParams, authorID string) ([]Post, error) {
	return s.repository.getPostsByAuthorID(ctx, params, authorID)
}
