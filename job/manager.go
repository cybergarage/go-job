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
	// RegisterJob registers a job in the registry. If a job with the same kind is already registered,
	// it will be overwritten with the new job.
	// If the registered job has a schedule configuration, it will be automatically scheduled and
	// the resulting job instance will be returned.
	// If the job has no schedule configuration, nil will be returned for the instance.
	RegisterJob(job Job) (Instance, error)
	// UnregisterJob removes a job from the registry by its kind.
	UnregisterJob(kind Kind) error
	// ListJobs returns a slice of all registered jobs.
	ListJobs() ([]Job, error)
	// LookupJob looks up a job by its kind in the registry.
	LookupJob(kind Kind) (Job, bool)
	// ScheduleJob schedules a job instance with the given job and options.
	// It creates a new job instance and enqueues it in the job queue.
	// If the schedule option is not set, the job instance will be scheduled to run immediately as default.
	ScheduleJob(job Job, opts ...any) (Instance, error)
	// ScheduleRegisteredJob schedules a registered job by its kind with the given options.
	// If the job is not registered, an error will be returned.
	// It creates a new job instance and enqueues it in the job queue.
	// If the schedule option is not set, the job instance will be scheduled to run immediately as default.
	ScheduleRegisteredJob(kind Kind, opts ...any) (Instance, error)
	// ListInstances returns all job instances which are currently scheduled, processing, completed, or terminated after the manager started.
	ListInstances() ([]Instance, error)
	// LookupInstances looks up all job instances which match the specified query.
	LookupInstances(query Query) ([]Instance, error)
	// LookupHistory retrieves all state records for a job instance, sorted by timestamp.
	LookupInstanceHistory(job Instance) (InstanceHistory, error)
	// LookupLogs retrieves all logs for a job instance.
	LookupInstanceLogs(job Instance) ([]Log, error)
	// Start starts the job manager.
	Start() error
	// Stop stops the job manager.
	Stop() error
	// Clear clears all jobs and history from the job manager without registered jobs.
	Clear() error
	// StopWithWait stops the job manager and waits for all jobs to complete.
	StopWithWait() error
	// ResizeWorkers scales the number of workers in the group.
	ResizeWorkers(num int) error
	// NumWorkers returns the number of workers in the group.
	NumWorkers() int
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

// RegisterJob registers a job in the registry. If a job with the same kind is already registered,
// it will be overwritten with the new job.
// If the registered job has a schedule configuration, it will be automatically scheduled and
// the resulting job instance will be returned.
// If the job has no schedule configuration, nil will be returned for the instance.
func (mgr *manager) RegisterJob(job Job) (Instance, error) {
	err := mgr.Repository.RegisterJob(job)
	if err != nil {
		return nil, fmt.Errorf("failed to register job: %w", err)
	}
	if !job.Schedule().IsScheduled() {
		return nil, nil
	}
	return mgr.ScheduleJob(job)
}

// ScheduleRegisteredJob schedules a registered job by its kind with the given options.
// If the job is not registered, an error will be returned.
// It creates a new job instance and enqueues it in the job queue.
// If the schedule option is not set, the job instance will be scheduled to run immediately as default.
func (mgr *manager) ScheduleRegisteredJob(kind Kind, opts ...any) (Instance, error) {
	job, ok := mgr.LookupJob(kind)
	if !ok {
		return nil, fmt.Errorf("registered job not found: %s", kind)
	}
	return mgr.ScheduleJob(job, opts...)
}

// ScheduleJob schedules a job instance with the given job and options.
// It creates a new job instance and enqueues it in the job queue.
// If the schedule option is not set, the job instance will be scheduled to run immediately as default.
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
	if err := ji.UpdateState(JobCreated, opts...); err != nil {
		return nil, err
	}

	if err := mgr.ScheduleJobInstance(ji); err != nil {
		return nil, err
	}
	if err := ji.UpdateState(JobScheduled, opts...); err != nil {
		return nil, err
	}

	return ji, nil
}

// ListInstance returns a list of all job instances which are currently scheduled, processing, completed, or terminated after the manager started.
func (mgr *manager) ListInstance() ([]Instance, error) {
	return mgr.LookupInstances(NewQuery())
}

// ListInstances returns all job instances which are currently scheduled, processing, completed, or terminated after the manager started.
func (mgr *manager) ListInstances() ([]Instance, error) {
	return mgr.LookupInstances(NewQuery())
}

// LookupInstances looks up all job instances which match the specified query.
func (mgr *manager) LookupInstances(query Query) ([]Instance, error) {
	mgr.Lock()
	defer mgr.Unlock()

	matchQuery := func(instance Instance, query Query) bool {
		if query == nil {
			return true // No query means match all
		}
		uuid, ok := query.UUID()
		if ok && (instance.UUID() != uuid) {
			return false
		}
		kind, ok := query.Kind()
		if ok && (instance.Kind() != kind) {
			return false
		}
		state, ok := query.State()
		if ok && (instance.State() != state) {
			return false
		}
		return true
	}

	var instances []Instance

	queueInstances, err := NewInstancesFromQueue(mgr.Queue())
	if err != nil {
		return nil, err
	}
	for _, instance := range queueInstances {
		if matchQuery(instance, query) {
			instances = append(instances, instance)
		}
	}

	history, err := mgr.ListHistory()
	if err != nil {
		return nil, err
	}
	historyInstances, err := NewInstancesFromHistory(history)
	if err != nil {
		return nil, err
	}
	for _, instance := range historyInstances {
		if matchQuery(instance, query) {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

// LookupHistory retrieves all state records for a job instance, sorted by timestamp.
func (mgr *manager) LookupInstanceHistory(job Instance) (InstanceHistory, error) {
	mgr.Lock()
	defer mgr.Unlock()

	return mgr.Repository.LookupHistory(job)
}

// LookupLogs retrieves all logs for a job instance.
func (mgr *manager) LookupInstanceLogs(job Instance) ([]Log, error) {
	mgr.Lock()
	defer mgr.Unlock()

	return mgr.Repository.LookupLogs(job)
}

// Start starts the job manager.
func (mgr *manager) Start() error {
	starters := []func() error{
		mgr.store.Start,
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
		mgr.store.Stop,
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
