package posts

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redisClient *redis.Client
}

func NewCache(redisClient *redis.Client) *Cache {
	return &Cache{redisClient: redisClient}
}

const postsCacheVersionKey = "posts:version"

func (c *Cache) generateKey(ctx context.Context, params listParams) (string, error) {
	version, err := c.redisClient.Get(ctx, postsCacheVersionKey).Result()
	if err == redis.Nil {
		version = "1"
	} else if err != nil {
		return "", err
	}

	return fmt.Sprintf("posts:v=%s:limit=%d:offset=%d", version, params.Limit, params.Offset), nil
}

func (c *Cache) getPosts(ctx context.Context, params listParams) ([]Post, error) {
	key, err := c.generateKey(ctx, params)
	if err != nil {
		return nil, err
	}

	val, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var posts []Post
	err = json.Unmarshal([]byte(val), &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (c *Cache) cachePosts(ctx context.Context, params listParams, posts []Post) error {
	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	key, err := c.generateKey(ctx, params)
	if err != nil {
		return err
	}

	err = c.redisClient.Set(ctx, key, data, time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) invalidatePosts(ctx context.Context) error {
	return c.redisClient.Incr(ctx, postsCacheVersionKey).Err()
}
