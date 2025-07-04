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
	"errors"
	"fmt"
	"sync"
)

// Manager is an interface that defines methods for managing jobs.
type Manager interface {
	// JobRegistry provides access to the job registry.
	JobRegistry
	// JobScheduler provides access to the job scheduler.
	Scheduler
	// ScheduleRegisteredJob schedules a registered job by its kind with the provided options.
	ScheduleRegisteredJob(kind Kind, opts ...any) error
	// Start starts the job manager.
	Start() error
	// Stop stops the job manager.
	Stop() error
}
type manager struct {
	sync.Mutex
	logger  Logger
	store   Store
	queue   Queue
	workers []Worker
	*scheduler
	JobRegistry
}

// ManagerOption is a function that configures a job manager.
type ManagerOption func(*manager)

// WithManagerStore sets the store for the job manager.
func WithNumWorkers(num int) ManagerOption {
	return func(m *manager) {
		m.workers = make([]Worker, num)
	}
}

// WithLogger sets the logger for the job manager.
func WithLogger(logger Logger) ManagerOption {
	return func(m *manager) {
		m.logger = logger
	}
}

// WithManagerQueue sets the queue for the job manager.
func WithStore(store Store) ManagerOption {
	return func(m *manager) {
		m.store = store
	}
}

// NewManager creates a new instance of the job manager.
func NewManager(opts ...ManagerOption) *manager {
	mgr := &manager{
		Mutex:       sync.Mutex{},
		logger:      NewNullLogger(),
		store:       NewMemStore(),
		scheduler:   nil,
		JobRegistry: NewJobRegistry(),
		queue:       nil,
		workers:     make([]Worker, 1),
	}
	for _, opt := range opts {
		opt(mgr)
	}

	mgr.queue = NewQueue(WithQueueStore(mgr.store))

	mgr.scheduler = NewScheduler(WithSchedulerQueue(mgr.queue))

	for i := 0; i < len(mgr.workers); i++ {
		mgr.workers[i] = NewWorker(WithWorkerQueue(mgr.queue))
	}

	return mgr
}

// ScheduleRegisteredJob schedules a registered job by its kind with the provided options.
func (mgr *manager) ScheduleRegisteredJob(kind Kind, opts ...any) error {
	job, ok := mgr.LookupJob(kind)
	if !ok {
		return fmt.Errorf("job not found: %s", kind)
	}
	return mgr.ScheduleJob(job, opts...)
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
