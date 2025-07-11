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
	// Repository defines the methods for managing job registrations.
	Repository
	// ScheduleJob schedules a job instance with the given job and options.
	ScheduleJob(job Job, opts ...any) (Instance, error)
	// ScheduleRegisteredJob schedules a registered job by its kind with the provided options.
	ScheduleRegisteredJob(kind Kind, opts ...any) (Instance, error)
	// Start starts the job manager.
	Start() error
	// Stop stops the job manager.
	Stop() error
}

type manager struct {
	sync.Mutex
	logger  Logger
	store   Store
	workers []Worker
	Repository
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
		Mutex:   sync.Mutex{},
		logger:  NewNullLogger(),
		store:   NewLocalStore(),
		workers: make([]Worker, 1),
	}

	for _, opt := range opts {
		opt(mgr)
	}

	mgr.Repository = NewRepository(
		WithRepositoryStore(mgr.store),
	)

	for i := 0; i < len(mgr.workers); i++ {
		mgr.workers[i] = NewWorker(WithWorkerQueue(mgr.Queue()))
	}

	return mgr
}

// ScheduleRegisteredJob schedules a registered job by its kind with the provided options.
func (mgr *manager) ScheduleRegisteredJob(kind Kind, opts ...any) (Instance, error) {
	job, ok := mgr.LookupJob(kind)
	if !ok {
		return nil, fmt.Errorf("registered job not found: %s", kind)
	}
	return mgr.ScheduleJob(job, opts...)
}

// ScheduleJob schedules a job instance with the given job and options.
// It creates a new job instance and enqueues it in the job queue.
func (mgr *manager) ScheduleJob(job Job, opts ...any) (Instance, error) {
	opts = append(opts,
		WithExecutor(job.Handler().Executor()),
		WithErrorHandler(job.Handler().ErrorHandler()),
		WithResponseHandler(job.Handler().ResponseHandler()),
		WithCrontabSpec(job.Schedule().CrontabSpec()),
	)
	ji, err := NewInstance(opts...)
	if err != nil {
		mgr.LogInstanceRecord(ji, JobError)
		return nil, err
	}
	mgr.LogInstanceRecord(ji, JobCreated)
	if err := mgr.ScheduleJobInstance(ji); err != nil {
		return nil, err
	}
	mgr.LogInstanceRecord(ji, JobScheduled)
	return ji, nil
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
