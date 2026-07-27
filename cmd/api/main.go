package main

import (
	"context"
	"fmt"
	"go-learn/internal/auth"
	"go-learn/internal/validation"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func connectDB() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return pool
}

func main() {
	loadEnv()

	router := gin.Default()

	pool := connectDB()
	defer pool.Close()

	validation.RegisterCustomValidators()

	authRepository := auth.NewRepository(pool)
	authService := auth.NewService(authRepository)
	authHandler := auth.NewHandler(authService)

	protected := router.Group("/")
	protected.Use(authHandler.RequireAuth)

	router.POST("/auth/sign-up", authHandler.SignUp)
	router.POST("/auth/sign-in", authHandler.SignIn)
	protected.GET("/auth/self", authHandler.GetSelf)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
}
