// Copyright (C) 2025 The go-job Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package job

import (
	"context"
	"sort"
	"sync"
	"time"
)

// Queue is an interface that defines methods for managing a job queue.
type Queue interface {
	// Enqueue adds a job to the queue.
	Enqueue(job Instance) error
	// Dequeue removes and returns a job from the queue.
	Dequeue() (Instance, error)
	// HasJobs checks if there are any jobs in the queue.
	HasJobs() (bool, error)
	// Size returns the number of jobs in the queue.
	Size() (int, error)
	// Empty checks if the queue is empty.
	Empty() (bool, error)
	// Lock acquires a lock for the queue.
	Lock() error
	// Unlock releases the lock for the queue.
	Unlock() error
}

type queue struct {
	sync.Mutex
	store Store
}

// QueueOption is a function that configures a job queue.
type QueueOption func(*queue)

// WithQueueStore sets the store for the job queue.
func WithQueueStore(store Store) QueueOption {
	return func(q *queue) {
		q.store = store
	}
}

// NewQueue creates a new instance of the job queue.
func NewQueue(opts ...QueueOption) Queue {
	queue := &queue{
		Mutex: sync.Mutex{},
		store: nil,
	}
	for _, opt := range opts {
		opt(queue)
	}
	return queue
}

// Lock acquires a lock for the queue.
func (q *queue) Lock() error {
	q.Mutex.Lock()
	return nil
}

// Unlock releases the lock for the queue.
func (q *queue) Unlock() error {
	q.Mutex.Unlock()
	return nil
}

// Enqueue adds a job to the queue.
func (q *queue) Enqueue(job Instance) error {
	q.Lock()
	defer q.Unlock()
	return q.store.EnqueueInstance(context.Background(), job)
}

// Dequeue removes and returns a job from the queue.
func (q *queue) Dequeue() (Instance, error) {
	for {
		q.Lock()
		ctx := context.Background()

		jobs, err := q.store.ListInstances(ctx)
		if err != nil {
			q.Unlock()
			return nil, err
		}

		sort.Slice(jobs, func(i, j int) bool {
			return jobs[i].Before(jobs[j])
		})

		now := time.Now()
		for _, job := range jobs {
			scheduledAt := job.ScheduledAt()
			if scheduledAt.Before(now) {
				err := q.store.RemoveInstance(ctx, job)
				if err == nil {
					q.Unlock()
					return job, nil
				}
			}
		}
		q.Unlock()
		time.Sleep(1 * time.Second)
	}
}

// Size returns the number of jobs in the queue.
func (q *queue) Size() (int, error) {
	q.Lock()
	defer q.Unlock()
	jobs, err := q.store.ListInstances(context.Background())
	if err != nil {
		return 0, err
	}
	return len(jobs), nil
}

// HasJobs checks if there are any jobs in the queue.
func (q *queue) HasJobs() (bool, error) {
	size, err := q.Size()
	return size > 0, err
}

// Empty checks if the queue is empty.
func (q *queue) Empty() (bool, error) {
	size, err := q.Size()
	return size == 0, err
}
