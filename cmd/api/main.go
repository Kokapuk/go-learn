package main

import (
	"context"
	"errors"
	"fmt"
	"go-learn/internal/auth"
	"go-learn/internal/jobs"
	"go-learn/internal/posts"
	"go-learn/internal/ratelimiter"
	"go-learn/internal/validation"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

func connectDB(ctx context.Context) *pgxpool.Pool {
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

func connectRedis(ctx context.Context) *redis.Client {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	err = redisClient.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	return redisClient
}

func connectPublisher() *jobs.Publisher {
	publisher, err := jobs.NewPublisher()
	if err != nil {
		panic(err)
	}

	return publisher
}

func runServer(ctx context.Context, router *gin.Engine) {
	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		log.Printf("API listening on %s", server.Addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	loadEnv()

	router := gin.Default()

	pool := connectDB(ctx)
	defer pool.Close()

	redisClient := connectRedis(ctx)
	defer redisClient.Close()

	publisher := connectPublisher()
	defer publisher.Close()

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
	postsService := posts.NewService(postsRepository, postsCache, publisher)
	postsHandler := posts.NewHandler(postsService)

	protected.POST("/posts", postsHandler.CreatePost)
	router.GET("/posts", postsHandler.GetPosts)
	protected.GET("/posts/mine", postsHandler.GetOwningPosts)

	runServer(ctx, router)
}
