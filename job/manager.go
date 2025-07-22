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
	"os"
	"sync"
	"time"

	logger "github.com/cybergarage/go-logger/log"
)

// Manager is an interface that defines methods for managing jobs.
type Manager interface {
	// Repository defines the methods for managing job registrations.
	Repository
	// WorkerGroup returns the worker group associated with the manager.
	WorkerGroup
	// ScheduleJob schedules a job instance with the given job and options.
	ScheduleJob(job Job, opts ...any) (Instance, error)
	// ScheduleRegisteredJob schedules a registered job by its kind with the provided options.
	ScheduleRegisteredJob(kind Kind, opts ...any) (Instance, error)
	// Start starts the job manager.
	Start() error
	// Stop stops the job manager.
	Stop() error
	// Clear clears all jobs and history from the job manager without registered jobs.
	Clear() error
	// StopWithWait stops the job manager and waits for all jobs to complete.
	StopWithWait() error
}

type manager struct {
	sync.Mutex
	store Store
	*workerGroup
	Repository
}

// ManagerOption is a function that configures a job manager.
type ManagerOption func(*manager)

// WithManagerQueue sets the queue for the job manager.
func WithStore(store Store) ManagerOption {
	return func(m *manager) {
		m.store = store
	}
}

// NewManager creates a new instance of the job manager.
func NewManager(opts ...any) (Manager, error) {
	return newManager(opts...)
}

// NewManager creates a new instance of the job manager.
func newManager(opts ...any) (*manager, error) {
	mgr := &manager{
		Mutex:       sync.Mutex{},
		store:       NewLocalStore(),
		workerGroup: newWorkerGroup(WithNumWorkers(DefaultWorkerNum)),
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case ManagerOption:
			opt(mgr)
		case WorkerGroupOption:
			opt(mgr.workerGroup)
		default:
			return nil, fmt.Errorf("invalid option type %T for job manager", opt)
		}
	}

	mgr.Repository = NewRepository(
		WithRepositoryStore(mgr.store),
	)
	WithWorkerGroupQueue(mgr.Repository.Queue())(mgr.workerGroup)

	return mgr, nil
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
	jobOpts := []any{
		WithExecutor(job.Handler().Executor()),
		WithErrorHandler(job.Handler().ErrorHandler()),
		WithResponseHandler(job.Handler().ResponseHandler()),
		WithCrontabSpec(job.Schedule().CrontabSpec()),
		WithInstanceHistory(mgr.Repository),
	}
	jobOpts = append(jobOpts, opts...)
	ji, err := NewInstance(jobOpts...)
	if err != nil {
		return nil, err
	}
	if err := mgr.ScheduleJobInstance(ji); err != nil {
		return nil, err
	}
	if err := ji.UpdateState(JobScheduled); err != nil {
		return nil, err
	}
	return ji, nil
}

// Start starts the job manager.
func (mgr *manager) Start() error {
	logger.Infof("%s/%s", ProductName, Version)

	starters := []func() error{
		mgr.workerGroup.Start,
	}
	var errs error
	for _, starter := range starters {
		if err := starter(); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	logger.Infof("%s (PID:%d) started", ProductName, os.Getpid())

	return errs
}

// Stop stops the job manager.
func (mgr *manager) Stop() error {
	stoppers := []func() error{
		mgr.workerGroup.Stop,
	}
	var errs error
	for _, stopper := range stoppers {
		if err := stopper(); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	logger.Infof("%s (PID:%d) terminated", ProductName, os.Getpid())

	return errs
}

// Clear clears all jobs and history from the job manager without registered jobs.
func (mgr *manager) Clear() error {
	mgr.Lock()
	defer mgr.Unlock()

	closer := []func() error{
		mgr.Repository.Clear,
	}
	for _, close := range closer {
		if err := close(); err != nil {
			return fmt.Errorf("failed to clear job manager: %w", err)
		}
	}
	return nil
}

// StopWithWait stops the job manager and waits for all jobs to complete.
func (mgr *manager) StopWithWait() error {
	for {
		if noJobs, _ := mgr.Queue().Empty(); noJobs {
			break
		}
		time.Sleep(100 * time.Millisecond) // Wait for queue to empty
	}

	mgr.Queue().Lock()
	defer mgr.Queue().Unlock()

	err := mgr.workerGroup.StopWithWait()
	if err != nil {
		return nil
	}

	return mgr.Stop()
}
