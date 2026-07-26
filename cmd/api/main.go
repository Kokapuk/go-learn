package main

import (
	"context"
	"fmt"
	"go-learn/internal/auth"
	"go-learn/internal/validation"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func connectDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return conn
}

func main() {
	loadEnv()

	router := gin.Default()

	conn := connectDB()
	defer conn.Close(context.Background())

	validation.RegisterCustomValidators()

	authRepository := auth.NewRepository(conn)
	authService := auth.NewService(authRepository)
	authHandler := auth.NewHandler(authService)

	router.POST("/auth/sign-up", authHandler.SignUp)
	// mux.HandleFunc("GET /users/{id}", usersHandler.GetUser)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	err = router.Run(fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
}
