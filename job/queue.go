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
)

// Queue is an interface that defines methods for managing a job queue.
type Queue interface {
	// Enqueue adds a job to the queue.
	Enqueue(job Instance) error
	// Dequeue removes and returns a job from the queue.
	Dequeue() (Instance, error)
}

type jobQueue struct {
	sync.Mutex
	store Store
}

// QueueOption is a function that configures a job queue.
type QueueOption func(*jobQueue)

// WithQueueStore sets the store for the job queue.
func WithQueueStore(store Store) QueueOption {
	return func(q *jobQueue) {
		q.store = store
	}
}

// NewQueue creates a new instance of the job queue.
func NewQueue(opts ...QueueOption) Queue {
	queue := &jobQueue{
		Mutex: sync.Mutex{},
		store: nil,
	}
	for _, opt := range opts {
		opt(queue)
	}
	return queue
}

// Enqueue adds a job to the queue.
func (q *jobQueue) Enqueue(job Instance) error {
	q.Lock()
	defer q.Unlock()
	return q.store.StoreJob(context.Background(), job)
}

// Dequeue removes and returns a job from the queue.
func (q *jobQueue) Dequeue() (Instance, error) {
	q.Lock()
	defer q.Unlock()
	ctx := context.Background()
	jobs, err := q.store.ListJobs(ctx)
	if err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		return nil, nil // No jobs to dequeue
	}
	job := jobs[0]
	if err := q.store.RemoveJob(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}
