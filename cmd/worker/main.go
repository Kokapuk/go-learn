package main

import (
	"context"
	"encoding/json"
	"errors"
	"go-learn/internal/jobs"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file detected")
	}
}

func runWorker(ctx context.Context) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"jobs.failed",
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"jobs", // name
		true,   // durability
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		amqp.Table{
			amqp.QueueTypeArg:           amqp.QueueTypeQuorum,
			"x-delivery-limit":          int32(3),
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": "jobs.failed",
		},
	)
	if err != nil {
		panic(err)
	}

	if err := ch.Qos(1, 0, false); err != nil {
		panic(err)
	}
	const consumerTag = "worker-post-created"
	msgs, err := ch.Consume(
		q.Name,      // queue
		consumerTag, // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for d := range msgs {
			var message jobs.PostCreatedMessage

			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("invalid job payload: %v", err)

				if err := d.Nack(false, false); err != nil {
					log.Printf("reject invalid message: %v", err)
				}
				continue
			}

			switch message.Type {
			case jobs.PostCreated:
				log.Printf(
					"processing post-created job: post=%s author=%s created_at=%s",
					message.PostID,
					message.AuthorID,
					message.CreatedAt.Format(time.RFC3339),
				)

				time.Sleep(time.Second * 5)
				success := rand.Float64() > 0.5
				log.Println("Job result:", success)

				if !success {
					log.Println("job failed; requeueing it")

					if err := d.Nack(false, true); err != nil {
						log.Printf("requeue failed job: %v", err)
					}
					continue
				} else {
					log.Println("Job finished successfully:", success)
				}

			default:
				log.Printf("unknown job type: %q", message.Type)

				if err := d.Nack(false, false); err != nil {
					log.Printf("reject unknown job: %v", err)
				}
				continue
			}

			if err := d.Ack(false); err != nil {
				log.Printf("acknowledge message: %v", err)
			}
		}
	}()

	log.Println("Worker is running")
	select {
	case <-ctx.Done():
	case connCloseErr := <-conn.NotifyClose(make(chan *amqp.Error)):
		log.Println("conn closed:", connCloseErr)
	}
	log.Println("Shutdown signal received")
	if err := ch.Cancel(consumerTag, false); err != nil && !errors.Is(err, amqp.ErrClosed) {
		log.Println("cancel consumer:", err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Worker stopped gracefully")
	case <-time.After(30 * time.Second):
		log.Println("Shutdown timeout exceeded, forcing exit")
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	loadEnv()
	runWorker(ctx)
}
