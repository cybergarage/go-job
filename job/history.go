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
	LogInstanceRecord(job Instance, state JobState) error
	InstanceRecords(job Instance) []InstanceRecord
	LastInstanceRecord(job Instance) InstanceRecord
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

// LogInstanceRecord logs a state change for a job instance.
func (h *history) LogInstanceRecord(job Instance, state JobState) error {
	record := newInstanceRecord(job.UUID(), state)
	return h.store.LogInstanceRecord(context.Background(), job, record)
}

// InstanceRecords retrieves all state records for a job instance, sorted by timestamp.
func (h *history) InstanceRecords(job Instance) []InstanceRecord {
	records, err := h.store.ListInstanceRecords(context.Background(), job)
	if err != nil {
		return nil
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp().Before(records[j].Timestamp())
	})
	return records
}

// LastInstanceRecord retrieves the most recent state record for a job instance.
func (h *history) LastInstanceRecord(job Instance) InstanceRecord {
	records := h.InstanceRecords(job)
	if len(records) == 0 {
		return nil
	}
	return records[len(records)-1]
}
