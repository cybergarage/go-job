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
	"fmt"
	"sync"
	"time"
)

type localStore struct {
	sync.Mutex

	jobs    sync.Map
	history []InstanceState
	logs    []Log
}

// NewLocalStore creates a new in-memory job store.
func NewLocalStore() Store {
	return &localStore{
		Mutex:   sync.Mutex{},
		jobs:    sync.Map{},
		history: []InstanceState{},
		logs:    []Log{},
	}
}

// Name returns the name of the store.
func (store *localStore) Name() string {
	return "local"
}

// EnqueueInstance stores a job instance in the store.
func (store *localStore) EnqueueInstance(ctx context.Context, job Instance) error {
	store.jobs.Store(job.UUID(), job)
	return nil
}

// DequeueNextInstance retrieves and removes the highest priority job instance from the store. If no job instance is available, it returns nil.
func (store *localStore) DequeueNextInstance(ctx context.Context) (Instance, error) {
	now := time.Now()
	var nextJob Instance
	store.jobs.Range(func(key, value interface{}) bool {
		if job, ok := value.(Instance); ok {
			scheduledAt := job.ScheduledAt()
			if scheduledAt.Before(now) {
				switch {
				case nextJob == nil:
					nextJob = job
				case job.Before(nextJob):
					nextJob = job
				}
			}
		}
		return true
	})
	if nextJob == nil {
		return nil, nil
	}
	store.jobs.Delete(nextJob.UUID())
	return nextJob, nil
}

// ListInstances lists all job instances in the store.
func (store *localStore) ListInstances(ctx context.Context) ([]Instance, error) {
	jobs := make([]Instance, 0)
	store.jobs.Range(func(key, value interface{}) bool {
		if job, ok := value.(Instance); ok {
			jobs = append(jobs, job)
		}
		return true
	})
	return jobs, nil
}

// ClearInstances clears all job instances in the store.
func (store *localStore) ClearInstances(ctx context.Context) error {
	store.jobs.Range(func(key, value any) bool {
		store.jobs.Delete(key)
		return true
	})
	return nil
}

// LogInstanceState adds a new state record for a job instance.
func (store *localStore) LogInstanceState(ctx context.Context, state InstanceState) error {
	store.Lock()
	defer store.Unlock()
	store.history = append(store.history, state)
	return nil
}

// LookupInstanceHistory lists all state records for a job instance that match the specified query. The returned history is sorted by their timestamp.
func (store *localStore) LookupInstanceHistory(ctx context.Context, query Query) (InstanceHistory, error) {
	store.Lock()
	defer store.Unlock()
	var records []InstanceState
	for _, record := range store.history {
		if query.Matches(record) {
			records = append(records, record)
		}
	}
	return records, nil
}

// ClearInstanceHistory clears all state records for a job instance that match the specified filter.
func (store *localStore) ClearInstanceHistory(ctx context.Context, filter Filter) error {
	store.Lock()
	defer store.Unlock()
	states := make([]InstanceState, 0)
	for _, state := range store.history {
		if filter.Matches(state) {
			continue
		}
		states = append(states, state)
	}
	store.history = states
	return nil
}

// Logf logs a formatted message at the specified log level.
func (store *localStore) Logf(ctx context.Context, job Instance, logLevel LogLevel, format string, args ...any) error {
	store.Lock()
	defer store.Unlock()
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
func (store *localStore) Infof(ctx context.Context, job Instance, format string, args ...any) error {
	return store.Logf(ctx, job, LogInfo, format, args...)
}

// Warnf logs a warning message for a job instance.
func (store *localStore) Warnf(ctx context.Context, job Instance, format string, args ...any) error {
	return store.Logf(ctx, job, LogWarn, format, args...)
}

// Errorf logs an error message for a job instance.
func (store *localStore) Errorf(ctx context.Context, job Instance, format string, args ...any) error {
	return store.Logf(ctx, job, LogError, format, args...)
}

// LookupInstanceLogs lists all log entries for a job instance that match the specified query. The returned logs are sorted by their timestamp.
func (store *localStore) LookupInstanceLogs(ctx context.Context, query Query) ([]Log, error) {
	store.Lock()
	defer store.Unlock()
	var logs []Log
	for _, log := range store.logs {
		if query.Matches(log) {
			logs = append(logs, log)
		}
	}
	return logs, nil
}

// ClearInstanceLogs clears all log entries for a job instance that match the specified filter.
func (store *localStore) ClearInstanceLogs(ctx context.Context, filter Filter) error {
	store.Lock()
	defer store.Unlock()
	logs := []Log{}
	for _, log := range store.logs {
		if filter.Matches(log) {
			continue
		}
		logs = append(logs, log)
	}
	store.logs = logs
	return nil
}

// Start starts the local store.
func (store *localStore) Start() error {
	// No specific start logic for local store
	return nil
}

// Stop stops the local store.
func (store *localStore) Stop() error {
	// No specific stop logic for local store
	return nil
}

// Clear removes all key-value objects from the store.
func (store *localStore) Clear() error {
	store.Lock()
	defer store.Unlock()
	store.history = []InstanceState{}
	store.logs = []Log{}
	return nil
}
