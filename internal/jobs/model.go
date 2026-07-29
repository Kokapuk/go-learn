package jobs

import (
	"context"
	"time"
)

type jobStatus string

const (
	statusQueued    jobStatus = "queued"
	statusRunning   jobStatus = "running"
	statusCompleted jobStatus = "completed"
	statusFailed    jobStatus = "failed"
	statusTimedOut  jobStatus = "timedOut"
	statusCancelled jobStatus = "cancelled"
)

type jobTask func(ctx context.Context) error

type job struct {
	ID        int       `json:"id"`
	Status    jobStatus `json:"status"`
	Task      jobTask   `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
