package main

import (
	"context"
	"fmt"
	"go-learn/internal/users"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s (%s)",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

func main() {
	loadEnv()

	mux := http.NewServeMux()

	conn := connectDB()
	defer conn.Close(context.Background())

	usersRepository := users.NewRepository(conn)
	usersService := users.NewService(usersRepository)
	usersHandler := users.NewHandler(usersService)

	mux.HandleFunc("POST /users", usersHandler.CreateUser)
	// mux.HandleFunc("GET /users/{id}", usersHandler.GetUser)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	log.Println("Listening at port", port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), logging(mux))
	if err != nil {
		panic(err)
	}
}
