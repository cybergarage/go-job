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
)

// History is an interface that defines methods for managing the history of job instance state changes.
type History interface {
	// StateHistory provides methods for managing the state history of job instances.
	StateHistory
	// LogHistory provides methods for logging messages related to job instances.
	LogHistory
}

// StateHistory is an interface that defines methods for managing the state history of job instances.
type StateHistory interface {
	// LogProcessState logs a state change for a job instance.
	LogProcessState(job Instance, state JobState, opts ...InstanceStateOption) error
	// LookupHistory lists all state records for a job instance that match the specified query. The returned history is sorted by their timestamp.
	LookupHistory(query Query) (InstanceHistory, error)
}

// LogHistory is an interface that defines methods for logging messages related to job instances.
type LogHistory interface {
	// Infof logs an informational message for a job instance.
	Infof(job Instance, format string, args ...any) error
	// Warnf logs a warning message for a job instance.
	Warnf(job Instance, format string, args ...any) error
	// Errorf logs an error message for a job instance.
	Errorf(job Instance, format string, args ...any) error
	// LookupInstanceLogs lists all log entries for a job instance that match the specified query. The returned logs are sorted by their timestamp.
	LookupLogs(query Query) ([]Log, error)
}

// HistoryOption is a function that configures the job history.
type HistoryOption func(*history)

// WithHistoryStore sets the store for the job history.
func WithHistoryStore(store HistoryStore) HistoryOption {
	return func(h *history) {
		h.store = store
	}
}

// history keeps track of the state changes of a job.
type history struct {
	store HistoryStore
}

// NewHistory creates a new instance of the job history.
func NewHistory(opts ...HistoryOption) History {
	return newHistory(opts...)
}

// newHistory creates a new job state history.
func newHistory(opts ...HistoryOption) *history {
	history := &history{
		store: NewLocalStore(),
	}
	for _, opt := range opts {
		opt(history)
	}
	return history
}

// LogProcessState logs a state change for a job instance.
func (history *history) LogProcessState(job Instance, state JobState, opts ...InstanceStateOption) error {
	opts = append(opts, WithStateKind(job.Kind()))
	opts = append(opts, WithStateUUID(job.UUID()))
	opts = append(opts, WithStateJobState(state))
	record := newInstanceState(opts...)
	return history.store.LogInstanceState(context.Background(), record)
}

// LookupHistory lists all state records for a job instance that match the specified query. The returned history is sorted by their timestamp.
func (history *history) LookupHistory(query Query) (InstanceHistory, error) {
	records, err := history.store.LookupInstanceHistory(context.Background(), query)
	if err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp().Before(records[j].Timestamp())
	})
	return records, nil
}

// Infof logs an informational message for a job instance.
func (history *history) Infof(job Instance, format string, args ...any) error {
	return history.store.Infof(context.Background(), job, format, args...)
}

// Warnf logs a warning message for a job instance.
func (history *history) Warnf(job Instance, format string, args ...any) error {
	return history.store.Warnf(context.Background(), job, format, args...)
}

// Errorf logs an error message for a job instance.
func (history *history) Errorf(job Instance, format string, args ...any) error {
	return history.store.Errorf(context.Background(), job, format, args...)
}

// LookupInstanceLogs lists all log entries for a job instance that match the specified query. The returned logs are sorted by their timestamp.
func (history *history) LookupLogs(query Query) ([]Log, error) {
	return history.store.LookupInstanceLogs(context.Background(), query)
}
