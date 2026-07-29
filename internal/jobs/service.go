package jobs

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repository *Repository
	worker     *Worker
}

func NewService(repository *Repository, worker *Worker) *Service {
	return &Service{repository: repository, worker: worker}
}

var errQueueFull = errors.New("Job queue is full")

func (s *Service) enqueueJob() (int, error) {
	if len(s.worker.queue) == cap(s.worker.queue) {
		return 0, errQueueFull
	}

	var task jobTask = func(ctx context.Context) error {
		select {
		case <-time.After(5 * time.Second):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	job := s.repository.enqueueJob(task)
	s.worker.queue <- job.ID
	return job.ID, nil
}

func (s *Service) getJobStatus(jobID int) (jobStatus, error) {
	jobStatus, err := s.repository.getStatus(jobID)
	if err != nil {
		return "", err
	}

	return jobStatus, nil
}
