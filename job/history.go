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
	// ProcessHistory retrieves all state records for a job instance, sorted by timestamp.
	ProcessHistory(job Instance) (InstanceHistory, error)
}

// LogHistory is an interface that defines methods for logging messages related to job instances.
type LogHistory interface {
	Infof(job Instance, format string, args ...any) error
	Warnf(job Instance, format string, args ...any) error
	Errorf(job Instance, format string, args ...any) error
	ProcessLogs(job Instance) ([]Log, error)
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
func (h *history) LogProcessState(job Instance, state JobState, opts ...InstanceStateOption) error {
	record := newInstanceState(job.UUID(), state, opts...)
	return h.store.LogInstanceState(context.Background(), job, record)
}

// ProcessHistory retrieves all state records for a job instance, sorted by timestamp.
func (h *history) ProcessHistory(job Instance) (InstanceHistory, error) {
	records, err := h.store.ListInstanceHistory(context.Background(), job)
	if err != nil {
		return nil, err
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp().Before(records[j].Timestamp())
	})
	return records, nil
}

// Infof logs an informational message for a job instance.
func (h *history) Infof(job Instance, format string, args ...any) error {
	return h.store.Infof(context.Background(), job, format, args...)
}

// Warnf logs a warning message for a job instance.
func (h *history) Warnf(job Instance, format string, args ...any) error {
	return h.store.Warnf(context.Background(), job, format, args...)
}

// Errorf logs an error message for a job instance.
func (h *history) Errorf(job Instance, format string, args ...any) error {
	return h.store.Errorf(context.Background(), job, format, args...)
}

// ProcessLogs retrieves all logs for a job instance.
func (h *history) ProcessLogs(job Instance) ([]Log, error) {
	return h.store.ListInstanceLogs(context.Background(), job)
}
