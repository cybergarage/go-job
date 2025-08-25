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
	"fmt"
	"time"

	logger "github.com/cybergarage/go-logger/log"
)

// Manager is an interface that defines methods for managing jobs.
type Manager interface {
	// Store returns the job store.
	Store() Store

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
	// If no schedule option is set, the job instance will be scheduled to run immediately by default.
	// If the specified job is not registered, the manager will register the job automatically.
	ScheduleJob(job Job, opts ...any) (Instance, error)
	// ScheduleRegisteredJob schedules a registered job by its kind with the given options.
	// If the job is not registered, an error will be returned.
	// It creates a new job instance and enqueues it in the job queue.
	// If the schedule option is not set, the job instance will be scheduled to run immediately as default.
	ScheduleRegisteredJob(kind Kind, opts ...any) (Instance, error)
	// EnqueueInstance enqueues a job instance in the job queue.
	EnqueueInstance(job Instance) error
	// DequeueNextInstance returns the next scheduled job instance and dequeues it from the job queue.
	DequeueNextInstance() (Instance, error)
	// LookupInstances looks up all job instances which match the specified query.
	LookupInstances(query Query) ([]Instance, error)
	// CancelInstances cancels all job instances which match the specified query.
	CancelInstances(query Query) ([]Instance, error)
	// ListInstances returns all job instances which are currently scheduled, processing, completed, or terminated after the manager started.
	ListInstances() ([]Instance, error)

	// LookupHistory retrieves all state records for a job instance, sorted by timestamp.
	LookupInstanceHistory(query Query) (InstanceHistory, error)
	// ClearInstanceHistory clears all state records for a job instance that match the specified filter.
	ClearInstanceHistory(filter Filter) error

	// LookupLogs retrieves all logs for a job instance.
	LookupInstanceLogs(query Query) ([]Log, error)
	// ClearInstanceLogs clears all log entries for a job instance that match the specified filter.
	ClearInstanceLogs(filter Filter) error

	// Workers returns a list of all workers in the group.
	Workers() []Worker
	// ResizeWorkers scales the number of workers in the group.
	ResizeWorkers(ctx context.Context, num int) error
	// NumWorkers returns the number of workers in the group.
	NumWorkers() int

	// Start starts the job manager.
	Start() error
	// Stop stops the job manager.
	Stop() error
	// Wait waits for all scheduled jobs to complete or terminate.
	Wait(ctx context.Context) error
	// Clear clears all jobs and history from the job manager without registered jobs.
	Clear() error
}

type manager struct {
	*workerGroup
	repository

	store Store
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
		store:       NewLocalStore(),
		workerGroup: newWorkerGroup(WithNumWorkers(DefaultWorkerNum)),
		repository:  nil,
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

	mgr.repository = newRepository(
		withRepositoryStore(mgr.store),
	)
	withWorkerGroupManager(mgr)(mgr.workerGroup)

	return mgr, nil
}

// Store returns the job store.
func (mgr *manager) Store() Store {
	return mgr.store
}

// RegisterJob registers a job in the registry. If a job with the same kind is already registered,
// it will be overwritten with the new job.
// If the registered job has a schedule configuration, it will be automatically scheduled and
// the resulting job instance will be returned.
// If the job has no schedule configuration, nil will be returned for the instance.
func (mgr *manager) RegisterJob(job Job) (Instance, error) {
	err := mgr.repository.RegisterJob(job)
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
// If no schedule option is set, the job instance will be scheduled to run immediately by default.
// If the specified job is not registered, the manager will register the job automatically.
func (mgr *manager) ScheduleJob(job Job, opts ...any) (Instance, error) {
	_, ok := mgr.LookupJob(job.Kind())
	if !ok {
		_, err := mgr.RegisterJob(job)
		if err != nil {
			return nil, err
		}
	}

	jobOpts := []any{
		WithJob(job),
		WithInstanceHistory(mgr.repository),
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

	mQueuedJobs.WithLabelValues(ji.Kind()).Inc()

	return ji, nil
}

// EnqueueInstance enqueues a job instance in the job queue.
func (mgr *manager) EnqueueInstance(job Instance) error {
	return mgr.Queue().Enqueue(context.Background(), job)
}

// DequeueNextInstance returns the next scheduled job instance and dequeues it from the job queue.
func (mgr *manager) DequeueNextInstance() (Instance, error) {
	ctx := context.Background()

	instance, err := mgr.Queue().Dequeue(ctx)
	if err != nil {
		return nil, err
	}

	// If the instance has a executor handler, it means it was dequeued from the local store.

	if instance.Executor() != nil {
		return instance, nil
	}

	// If the instance has no handler, it means it was dequeued from a remote store.
	// In this case, we need to recreate the instance with the corresponding job information.

	job, ok := mgr.LookupJob(instance.Kind())
	if !ok {
		// Jobs are registered per manager, so if the job is not registered in this manager, we need to re-enqueue the instance that cannot be handled by the current manager.
		logger.Infof("manager does not have job registered for instance: %s", instance.Kind())
		logger.Infof("manager re-enqueueing instance: %s", instance.UUID())
		err := mgr.EnqueueInstance(instance) // Re-enqueue the instance
		if err != nil {
			logger.Errorf("failed to re-enqueue instance: %s", err)
		}
		return nil, fmt.Errorf("job not registered for instance: %s", instance.Kind())
	}

	// Recreate the instance with the corresponding job information, including the handler's executor.
	newInstance, err := NewInstance(
		WithJob(job),
		WithUUID(instance.UUID()),
		WithCreatedAt(instance.CreatedAt()),
		WithState(instance.State()),
		WithArguments(instance.Arguments()...),
		withInstanceStore(mgr.store),
	)
	if err != nil {
		return nil, err
	}
	return newInstance, nil
}

// LookupInstances looks up all job instances which match the specified query.
func (mgr *manager) LookupInstances(query Query) ([]Instance, error) {
	matchedInstances := []Instance{}

	allQueueInstances, err := newInstancesFromQueue(mgr.Queue())
	if err != nil {
		return nil, err
	}
	for _, queueInstance := range allQueueInstances {
		if !query.Matches(queueInstance) {
			continue
		}
		matchedInstances = append(matchedInstances, queueInstance)
	}

	history, err := mgr.LookupHistory(query)
	if err != nil {
		return nil, err
	}
	historyInstances, err := newInstancesFromHistory(history)
	if err != nil {
		return nil, err
	}
	for _, historyInstance := range historyInstances {
		if !query.Matches(historyInstance) {
			continue
		}
		matchedInstances = append(matchedInstances, historyInstance)
	}

	return matchedInstances, nil
}

// CancelInstances cancels all job instances which match the specified query.
func (mgr *manager) CancelInstances(query Query) ([]Instance, error) {
	canceledInstances := []Instance{}

	allQueueInstances, err := newInstancesFromQueue(mgr.Queue())
	if err != nil {
		return nil, err
	}
	for _, queueInstance := range allQueueInstances {
		if !query.Matches(queueInstance) {
			continue
		}
		if err := mgr.Queue().Remove(context.Background(), queueInstance); err != nil {
			return canceledInstances, err
		}
		if err := queueInstance.UpdateState(JobCanceled); err != nil {
			return canceledInstances, err
		}
		queueInstance.HandleTerminated(queueInstance, context.Canceled)
		canceledInstances = append(canceledInstances, queueInstance)
	}

	for _, worker := range mgr.Workers() {
		workerInstance, ok := worker.ProcessingInstance()
		if !ok {
			continue
		}
		if err := worker.Cancel(); err != nil {
			return canceledInstances, err
		}
		canceledInstances = append(canceledInstances, workerInstance)
	}

	return canceledInstances, nil
}

// ListInstances returns all job instances which are currently scheduled, processing, completed, or terminated after the manager started.
func (mgr *manager) ListInstances() ([]Instance, error) {
	return mgr.LookupInstances(NewQuery())
}

// LookupHistory retrieves all state records for a job instance, sorted by timestamp.
func (mgr *manager) LookupInstanceHistory(query Query) (InstanceHistory, error) {
	return mgr.LookupHistory(query)
}

// ClearInstanceHistory clears all state records for a job instance that match the specified filter.
func (mgr *manager) ClearInstanceHistory(filter Filter) error {
	return mgr.repository.ClearHistory(filter)
}

// LookupLogs retrieves all logs for a job instance.
func (mgr *manager) LookupInstanceLogs(query Query) ([]Log, error) {
	return mgr.LookupLogs(query)
}

// ClearInstanceLogs clears all log entries for a job instance that match the specified filter.
func (mgr *manager) ClearInstanceLogs(filter Filter) error {
	return mgr.repository.ClearLogs(filter)
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
	return errs
}

// Clear clears all jobs and history from the job manager without registered jobs.
func (mgr *manager) Clear() error {
	cleaners := []func() error{
		mgr.store.Clear,
		mgr.repository.Clear,
	}
	for _, cleaner := range cleaners {
		if err := cleaner(); err != nil {
			return fmt.Errorf("failed to clear job manager: %w", err)
		}
	}
	return nil
}

// Wait waits for all scheduled jobs to complete or terminate.
func (mgr *manager) Wait(ctx context.Context) error {
	for {
		if noJobs, _ := mgr.Queue().Empty(ctx); noJobs {
			break
		}
		if deadline, ok := ctx.Deadline(); ok && !deadline.IsZero() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(1 * time.Second):
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	err := mgr.workerGroup.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}
