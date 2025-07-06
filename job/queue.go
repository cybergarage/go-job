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
	"sync"
	"time"
)

// Queue is an interface that defines methods for managing a job queue.
type Queue interface {
	// Enqueue adds a job to the queue.
	Enqueue(job Instance) error
	// Dequeue removes and returns a job from the queue.
	Dequeue() (Instance, error)
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

// Enqueue adds a job to the queue.
func (q *queue) Enqueue(job Instance) error {
	q.Lock()
	defer q.Unlock()
	return q.store.StoreJob(context.Background(), job)
}

// Dequeue removes and returns a job from the queue.
func (q *queue) Dequeue() (Instance, error) {
	for {
		q.Lock()
		ctx := context.Background()
		jobs, err := q.store.ListJobs(ctx)
		if err != nil {
			q.Unlock()
			return nil, err
		}
		for _, job := range jobs {
			if job.ScheduledAt().After(time.Now()) {
				err := q.store.RemoveJob(ctx, job)
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
