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
	"errors"
	"sync"
)

type manager struct {
	sync.Mutex
	store   Store
	queue   Queue
	workers []Worker
}

// ManagerOption is a function that configures a job manager.
type ManagerOption func(*manager)

// WithManagerStore sets the store for the job manager.
func WithManagerNumWorkers(num int) ManagerOption {
	return func(m *manager) {
		m.workers = make([]Worker, num)
		for i := range num {
			m.workers[i] = NewWorker()
		}
	}
}

// NewManager creates a new instance of the job manager.
func NewManager(opts ...ManagerOption) *manager {
	mgr := &manager{
		store:   NewMemStore(),
		queue:   nil,
		workers: make([]Worker, 0),
	}
	for _, opt := range opts {
		opt(mgr)
	}
	mgr.queue = NewQueue(WithQueueStore(mgr.store))
	return mgr
}

// ProcessJob processes a job by executing it and updating its state.
func (mgr *manager) ProcessJob(ctx context.Context, job Job) error {
	if job == nil {
		return nil
	}
	// Implement the logic to process the job
	// For example, you might want to execute the job's command or function
	// and update its state in the database or in-memory store.
	return nil
}

// ScheduleJob schedules a job to run at a specific time or interval.
func (mgr *manager) ScheduleJob(ctx context.Context, job Job) error {
	// Implement the logic to schedule the job
	// For example, you might want to add the job to a queue or a scheduler
	// that will execute it at the specified time or interval.
	return nil
}

// CancelJob cancels a scheduled job.
func (mgr *manager) CancelJob(ctx context.Context, job Job) error {
	// Implement the logic to cancel the job
	// For example, you might want to remove the job from the queue or scheduler
	// and update its state to "canceled" in the database or in-memory store.
	return nil
}

// ListJobs lists all jobs with the specified state.
func (mgr *manager) ListJobs(ctx context.Context, state JobState) ([]Job, error) {
	storeJobs, err := mgr.store.ListJobs(ctx)
	if err != nil {
		return nil, err
	}
	jobs := make([]Job, 0)
	for _, job := range storeJobs {
		if job.State() == state {
			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}

// Start starts the job manager.
func (mgr *manager) Start() error {
	for _, w := range mgr.workers {
		if err := w.Start(); err != nil {
			return errors.Join(err, mgr.Stop())
		}
	}
	return nil
}

// Stop stops the job manager.
func (mgr *manager) Stop() error {
	for _, w := range mgr.workers {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	return nil
}

// ScaleWorkers scales the number of workers for the job manager.
func (mgr *manager) ScaleWorkers(num int) error {
	if num < 0 {
		return errors.New("number of workers cannot be negative")
	}
	if num == len(mgr.workers) {
		return nil
	}

	if !mgr.TryLock() {
		return errors.New("manager is scaling workers")
	}
	defer mgr.Unlock()

	if num > len(mgr.workers) {
		for i := len(mgr.workers); i < num; i++ {
			worker := NewWorker()
			if err := worker.Start(); err != nil {
				return err
			}
			mgr.workers = append(mgr.workers, worker)
		}
	} else {
		for i := num; i < len(mgr.workers); i++ {
			if err := mgr.workers[i].Stop(); err != nil {
				return err
			}
		}
		mgr.workers = mgr.workers[:num]
	}
	return nil
}
