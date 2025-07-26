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

package store

import (
	"context"
	"errors"
	"time"

	"github.com/cybergarage/go-job/job"
	"github.com/cybergarage/go-job/job/plugins/store/kv"
)

type kvStore struct {
	kv.Store
}

// NewKVStore creates a new key-value store instance.
func NewStoreWithKvStore(store kv.Store) *kvStore /*job.Store*/ {
	return &kvStore{
		Store: store,
	}
}

// EnqueueInstance stores a job instance in the store.
func (store *kvStore) EnqueueInstance(ctx context.Context, job job.Instance) error {
	obj, err := kv.NewObjectFromInstance(job)
	if err != nil {
		return err
	}
	tx, err := store.Transact(ctx, true)
	if err != nil {
		return err
	}
	if err := tx.Set(ctx, obj); err != nil {
		return errors.Join(err, tx.Cancel(ctx))
	}
	return tx.Commit(ctx)
}

// DequeueNextInstance retrieves and removes the highest priority job instance from the store. If no job instance is available, it returns nil.
func (store *kvStore) DequeueNextInstance(ctx context.Context) (job.Instance, error) {
	now := time.Now()

	tx, err := store.Transact(ctx, true)
	if err != nil {
		return nil, err
	}

	rs, err := tx.GetRange(ctx, kv.NewInstanceListKey())
	if err != nil {
		return nil, err
	}

	var nextJob job.Instance
	for rs.Next() {
		obj, err := rs.Object()
		if err != nil {
			return nil, errors.Join(err, tx.Cancel(ctx))
		}
		job, err := kv.NewInstanceFromBytes(obj.Bytes())
		if err != nil {
			return nil, errors.Join(err, tx.Cancel(ctx))
		}
		scheduledAt := job.ScheduledAt()
		if scheduledAt.After(now) {
			continue
		}
		switch {
		case nextJob == nil:
			nextJob = job
		case job.Before(nextJob):
			nextJob = job
		}
	}

	if nextJob == nil {
		return nil, tx.Commit(ctx)
	}

	key := kv.NewInstanceKeyFrom(nextJob)

	err = tx.Remove(ctx, key)
	if err != nil {
		return nil, errors.Join(err, tx.Cancel(ctx))
	}

	return nextJob, tx.Commit(ctx)
}

// DequeueInstance removes a job instance from the store by its unique identifier.
func (store *kvStore) DequeueInstance(ctx context.Context, job job.Instance) error {
	tx, err := store.Transact(ctx, true)
	if err != nil {
		return err
	}

	key := kv.NewInstanceKeyFrom(job)

	err = tx.Remove(ctx, key)
	if err != nil {
		return errors.Join(err, tx.Cancel(ctx))
	}

	return tx.Commit(ctx)
}

// ListInstances lists all job instances in the store.
func (store *kvStore) ListInstances(ctx context.Context) ([]job.Instance, error) {
	tx, err := store.Transact(ctx, true)
	if err != nil {
		return nil, err
	}

	rs, err := tx.GetRange(ctx, kv.NewInstanceListKey())
	if err != nil {
		return nil, err
	}

	jobs := make([]job.Instance, 0)
	for rs.Next() {
		obj, err := rs.Object()
		if err != nil {
			return nil, err
		}
		job, err := kv.NewInstanceFromBytes(obj.Bytes())
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// ClearInstances clears all job instances in the store.
func (store *kvStore) ClearInstances(ctx context.Context) error {
	tx, err := store.Transact(ctx, true)
	if err != nil {
		return err
	}
	err = tx.RemoveRange(ctx, kv.NewInstanceListKey())
	if err != nil {
		return errors.Join(err, tx.Cancel(ctx))
	}
	return tx.Commit(ctx)
}

/*
// LogInstanceState adds a new state record for a job instance.
func (store *kvStore) LogInstanceState(ctx context.Context, job job.Instance, state InstanceState) error {
	if job == nil {
		return nil
	}
	store.history = append(store.history, state)
	return nil
}

// LookupInstanceHistory lists all state records for a job instance.
func (store *kvStore) LookupInstanceHistory(ctx context.Context, job job.Instance) (InstanceHistory, error) {
	if job == nil {
		return nil, nil
	}
	var records []InstanceState
	for _, record := range store.history {
		if record.UUID() == job.UUID() {
			records = append(records, record)
		}
	}
	return records, nil
}

// ListInstanceHistory lists all state records for all job instances.
func (store *kvStore) ListInstanceHistory(ctx context.Context) (InstanceHistory, error) {
	if len(store.history) == 0 {
		return nil, nil
	}
	return store.history, nil
}

// ClearInstanceHistory clears all state records for a job instance.
func (store *kvStore) ClearInstanceHistory(ctx context.Context) error {
	store.history = []InstanceState{}
	return nil
}

// Logf logs a formatted message at the specified log level.
func (store *kvStore) Logf(ctx context.Context, job job.Instance, logLevel LogLevel, format string, args ...any) error {
	log := NewLog(
		WithLogKind(job.Kind()),
		WithLogUUID(job.UUID()),
		WithLogLevel(logLevel),
		WithLogMessage(fmt.Sprintf(format, args...)),
	)
	store.logs = append(store.logs, log)
	return nil
}

// Infof logs an informational message for a job instance.
func (store *kvStore) Infof(ctx context.Context, job job.Instance, format string, args ...any) error {
	return store.Logf(ctx, job, LogInfo, format, args...)
}

// Warnf logs a warning message for a job instance.
func (store *kvStore) Warnf(ctx context.Context, job job.Instance, format string, args ...any) error {
	return store.Logf(ctx, job, LogWarn, format, args...)
}

// Errorf logs an error message for a job instance.
func (store *kvStore) Errorf(ctx context.Context, job job.Instance, format string, args ...any) error {
	return store.Logf(ctx, job, LogError, format, args...)
}

// LookupInstanceLogs lists all log entries for a job instance.
func (store *kvStore) LookupInstanceLogs(ctx context.Context, job job.Instance) ([]Log, error) {
	var logs []Log
	for _, log := range store.logs {
		if log.UUID() == job.UUID() {
			logs = append(logs, log)
		}
	}
	return logs, nil
}

// ClearInstanceLogs clears all log entries for a job instance.
func (store *kvStore) ClearInstanceLogs(ctx context.Context) error {
	store.logs = []Log{}
	return nil
}
*/
