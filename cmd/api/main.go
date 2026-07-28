package main

import (
	"context"
	"fmt"
	"go-learn/internal/auth"
	"go-learn/internal/posts"
	"go-learn/internal/ratelimiter"
	"go-learn/internal/validation"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file detected")
	}
}

func connectDB() *pgxpool.Pool {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		panic(err)
	}

	return pool
}

func connectRedis() *redis.Client {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return redisClient
}

func main() {
	loadEnv()

	router := gin.Default()

	pool := connectDB()
	defer pool.Close()

	redisClient := connectRedis()
	defer redisClient.Close()

	limiter := ratelimiter.NewRateLimiter(redisClient)

	validation.RegisterCustomValidators()

	authRepository := auth.NewRepository(pool)
	authService := auth.NewService(authRepository)
	authHandler := auth.NewHandler(authService)

	protected := router.Group("/")
	protected.Use(authHandler.RequireAuth)

	router.POST("/auth/sign-up", limiter.Limit("auth:sign-up", 3, time.Hour), authHandler.SignUp)
	router.POST("/auth/sign-in", limiter.Limit("auth:sign-in", 5, time.Minute), authHandler.SignIn)
	protected.GET("/auth/self", authHandler.GetSelf)

	postsRepository := posts.NewRepository(pool)
	postsCache := posts.NewCache(redisClient)
	postsService := posts.NewService(postsRepository, postsCache)
	postsHandler := posts.NewHandler(postsService)

	protected.POST("/posts", postsHandler.CreatePost)
	router.GET("/posts", postsHandler.GetPosts)
	protected.GET("/posts/mine", postsHandler.GetOwningPosts)

	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
}
