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
)

type localStore struct {
	jobs    sync.Map
	history []InstanceState
	logs    []Log
}

// NewLocalStore creates a new in-memory job store.
func NewLocalStore() Store {
	return &localStore{
		jobs:    sync.Map{},
		history: []InstanceState{},
		logs:    []Log{},
	}
}

// Name returns the name of the store.
func (s *localStore) Name() string {
	return "local"
}

// EnqueueInstance stores a job instance in the store.
func (s *localStore) EnqueueInstance(ctx context.Context, job Instance) error {
	if job == nil {
		return nil
	}
	s.jobs.Store(job.UUID(), job)
	return nil
}

// DequeueInstance removes a job instance from the store by its unique identifier.
func (s *localStore) DequeueInstance(ctx context.Context, job Instance) error {
	if job == nil {
		return nil
	}
	s.jobs.Delete(job.UUID())
	return nil
}

// ListInstances lists all job instances in the store.
func (s *localStore) ListInstances(ctx context.Context) ([]Instance, error) {
	jobs := make([]Instance, 0)
	s.jobs.Range(func(key, value interface{}) bool {
		if job, ok := value.(Instance); ok {
			jobs = append(jobs, job)
		}
		return true
	})
	return jobs, nil
}

// ClearInstances clears all job instances in the store.
func (s *localStore) ClearInstances(ctx context.Context) error {
	s.jobs.Range(func(key, value any) bool {
		s.jobs.Delete(key)
		return true
	})
	return nil
}

// LogInstanceState adds a new state record for a job instance.
func (s *localStore) LogInstanceState(ctx context.Context, job Instance, state InstanceState) error {
	if job == nil {
		return nil
	}
	s.history = append(s.history, state)
	return nil
}

// LookupInstanceHistory lists all state records for a job instance.
func (s *localStore) LookupInstanceHistory(ctx context.Context, job Instance) (InstanceHistory, error) {
	if job == nil {
		return nil, nil
	}
	var records []InstanceState
	for _, record := range s.history {
		if record.UUID() == job.UUID() {
			records = append(records, record)
		}
	}
	return records, nil
}

// ListInstanceHistory lists all state records for all job instances.
func (s *localStore) ListInstanceHistory(ctx context.Context) (InstanceHistory, error) {
	if len(s.history) == 0 {
		return nil, nil
	}
	return s.history, nil
}

// ClearInstanceHistory clears all state records for a job instance.
func (s *localStore) ClearInstanceHistory(ctx context.Context) error {
	s.history = []InstanceState{}
	return nil
}

// Logf logs a formatted message at the specified log level.
func (s *localStore) Logf(ctx context.Context, job Instance, logLevel LogLevel, format string, args ...any) error {
	log := NewLog(
		WithLogKind(job.Kind()),
		WithLogUUID(job.UUID()),
		WithLogLevel(logLevel),
		WithLogMessage(fmt.Sprintf(format, args...)),
	)
	s.logs = append(s.logs, log)
	return nil
}

// Infof logs an informational message for a job instance.
func (s *localStore) Infof(ctx context.Context, job Instance, format string, args ...any) error {
	return s.Logf(ctx, job, LogInfo, format, args...)
}

// Warnf logs a warning message for a job instance.
func (s *localStore) Warnf(ctx context.Context, job Instance, format string, args ...any) error {
	return s.Logf(ctx, job, LogWarn, format, args...)
}

// Errorf logs an error message for a job instance.
func (s *localStore) Errorf(ctx context.Context, job Instance, format string, args ...any) error {
	return s.Logf(ctx, job, LogError, format, args...)
}

// LookupInstanceLogs lists all log entries for a job instance.
func (s *localStore) LookupInstanceLogs(ctx context.Context, job Instance) ([]Log, error) {
	var logs []Log
	for _, log := range s.logs {
		if log.UUID() == job.UUID() {
			logs = append(logs, log)
		}
	}
	return logs, nil
}

// ClearInstanceLogs clears all log entries for a job instance.
func (s *localStore) ClearInstanceLogs(ctx context.Context) error {
	s.logs = []Log{}
	return nil
}
