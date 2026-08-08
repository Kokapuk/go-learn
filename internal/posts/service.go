package posts

import (
	"context"
	"errors"
	"go-learn/internal/jobs"
	"log"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repository *Repository
	cache      *Cache
	publisher  *jobs.Publisher
}

func NewService(repository *Repository, cache *Cache, publisher *jobs.Publisher) *Service {
	return &Service{repository: repository, cache: cache, publisher: publisher}
}

func (s *Service) createPost(ctx context.Context, params createPostParams) (Post, error) {
	post, err := s.repository.createPost(ctx, params)
	if err != nil {
		return Post{}, err
	}

	if err := s.cache.invalidatePosts(ctx); err != nil {
		log.Println(err)
	}

	message := jobs.PostCreatedMessage{
		Type:      jobs.PostCreated,
		PostID:    post.ID,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt,
	}
	if err := s.publisher.PublishPostCreated(ctx, message); err != nil {
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
