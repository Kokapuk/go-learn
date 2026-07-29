package jobs

import (
	"context"
	"errors"
	"log"
	"time"
)

type Worker struct {
	queue      chan int
	repository *Repository
}

func NewWorker(repository *Repository) *Worker {
	return &Worker{
		queue:      make(chan int, 100),
		repository: repository,
	}
}

func (w *Worker) Run(workerCtx context.Context) {
	for {
		select {
		case <-workerCtx.Done():
			return

		case jobID := <-w.queue:
			err := w.repository.setStatus(jobID, statusRunning)
			if err != nil {
				log.Println(err)
				continue
			}

			task, err := w.repository.getTask(jobID)
			if err != nil {
				log.Println(err)

				err := w.repository.setStatus(jobID, statusFailed)
				if err != nil {
					log.Println(err)
				}

				continue
			}

			jobCtx, cancel := context.WithTimeout(workerCtx, 10*time.Second)
			err = task(jobCtx)
			cancel()

			switch {
			case err == nil:
				err := w.repository.setStatus(jobID, statusCompleted)
				if err != nil {
					log.Println(err)
				}
			case errors.Is(err, context.DeadlineExceeded):
				err := w.repository.setStatus(jobID, statusTimedOut)
				if err != nil {
					log.Println(err)
				}
			case errors.Is(err, context.Canceled):
				err := w.repository.setStatus(jobID, statusCancelled)
				if err != nil {
					log.Println(err)
				}
			default:
				log.Println(err)

				err := w.repository.setStatus(jobID, statusFailed)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
