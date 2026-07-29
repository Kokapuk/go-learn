package jobs

import (
	"errors"
	"sync"
	"time"
)

type Repository struct {
	mu     sync.RWMutex
	nextId int
	jobs   map[int]*job
}

func NewRepository() *Repository {
	return &Repository{nextId: 1, jobs: make(map[int]*job)}
}

var errJobDoesNotExist = errors.New("Job with that id does not exists")

func (r *Repository) enqueueJob(task jobTask) *job {
	r.mu.Lock()
	defer r.mu.Unlock()

	job := &job{ID: r.nextId, Status: statusQueued, Task: task, CreatedAt: time.Now()}
	r.jobs[r.nextId] = job
	r.nextId++
	return job
}

func (r *Repository) setStatus(id int, status jobStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[id]
	if !exists {
		return errJobDoesNotExist
	}

	job.Status = status
	return nil
}

func (r *Repository) getStatus(id int) (jobStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[id]
	if !exists {
		return "", errJobDoesNotExist
	}

	return job.Status, nil
}

func (r *Repository) getTask(id int) (jobTask, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[id]
	if !exists {
		return nil, errJobDoesNotExist
	}

	return job.Task, nil
}
