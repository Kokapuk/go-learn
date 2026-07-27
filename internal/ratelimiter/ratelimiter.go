package ratelimiter

import (
	"fmt"
	"go-learn/internal/httputil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient *redis.Client
}

func NewRateLimiter(redisClient *redis.Client) *RateLimiter {
	return &RateLimiter{redisClient: redisClient}
}

func (h *RateLimiter) Limit(key string, limit int, ttl time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate-limit:%v:%v", key, c.ClientIP())

		pipe := h.redisClient.TxPipeline()
		incr := pipe.Incr(c, key)
		pipe.ExpireNX(c, key, ttl)
		if _, err := pipe.Exec(c); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, httputil.ErrorResponse{Message: "Something went wrong"})
			return
		}

		count := incr.Val()

		if count > int64(limit) {
			remainingTTL, _ := h.redisClient.TTL(c, key).Result()
			c.Header("Retry-After", strconv.Itoa(int(math.Ceil(remainingTTL.Seconds()))))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, httputil.ErrorResponse{Message: "Too many requests"})
			return
		}

		c.Next()
	}
}
