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
	"time"
)

// Queue is an interface that defines methods for managing a job queue.
type Queue interface {
	// Enqueue adds a job to the queue.
	Enqueue(ctx context.Context, job Instance) error
	// Dequeue removes and returns a job from the queue.
	Dequeue(ctx context.Context) (Instance, error)
	// List returns a list of all jobs in the queue.
	List(ctx context.Context) ([]Instance, error)
	// Size returns the number of jobs in the queue.
	Size(ctx context.Context) (int, error)
	// Empty checks if the queue is empty.
	Empty(ctx context.Context) (bool, error)
	// Clear clears all jobs in the queue.
	Clear(ctx context.Context) error
}

type queue struct {
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
		store: nil,
	}
	for _, opt := range opts {
		opt(queue)
	}
	return queue
}

// Enqueue adds a job to the queue.
func (q *queue) Enqueue(ctx context.Context, job Instance) error {
	return q.store.EnqueueInstance(ctx, job)
}

// Dequeue removes and returns a job from the queue.
func (q *queue) Dequeue(ctx context.Context) (Instance, error) {
	for {
		job, err := q.store.DequeueNextInstance(ctx)
		if err != nil {
			return nil, err
		}
		if job != nil {
			return job, nil
		}
		time.Sleep(1 * time.Second)
	}
}

// List returns a list of all jobs in the queue.
func (q *queue) List(ctx context.Context) ([]Instance, error) {
	jobs, err := q.store.ListInstances(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Before(jobs[j])
	})
	return jobs, nil
}

// Size returns the number of jobs in the queue.
func (q *queue) Size(ctx context.Context) (int, error) {
	jobs, err := q.store.ListInstances(ctx)
	if err != nil {
		return 0, err
	}
	return len(jobs), nil
}

// Empty checks if the queue is empty.
func (q *queue) Empty(ctx context.Context) (bool, error) {
	size, err := q.Size(ctx)
	return size == 0, err
}

// Clear clears all jobs in the queue.
func (q *queue) Clear(ctx context.Context) error {
	return q.store.ClearInstances(ctx)
}
