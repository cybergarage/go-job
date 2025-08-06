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
)

// Store defines the interface for job queue, history, and logging.
type Store interface {
	// Name returns the name of the store.
	Name() string
	// PendingStore provides methods for managing job instances.
	QueueStore
	// HistoryStore provides methods for managing job instance state history.
	HistoryStore
	// Start starts the store.
	Start() error
	// Stop stops the store.
	Stop() error
}

// QueueStore is an interface that defines methods for managing job instances in a pending state.
type QueueStore interface {
	// EnqueueInstance stores a job instance in the store.
	EnqueueInstance(ctx context.Context, job Instance) error
	// DequeueNextInstance retrieves and removes the highest priority job instance from the store. If no job instance is available, it returns nil.
	DequeueNextInstance(ctx context.Context) (Instance, error)
	// DequeueInstance removes a job instance from the store by its unique identifier.
	DequeueInstance(ctx context.Context, job Instance) error
	// ListInstances lists all job instances in the store.
	ListInstances(ctx context.Context) ([]Instance, error)
	// ClearInstances clears all job instances in the store.
	ClearInstances(ctx context.Context) error
}

// HistoryStore is an interface that defines methods for managing job instance state history.
type HistoryStore interface {
	// StateStore provides methods for managing job instance state history.
	StateStore
	// LogStore provides methods for logging job instance messages.
	LogStore
}

// StateStore is an interface that defines methods for managing job instance state history.
type StateStore interface {
	// LogInstanceState adds a new state record for a job instance.
	LogInstanceState(ctx context.Context, state InstanceState) error
	// LookupInstanceHistory lists all state records for a job instance. The returned history is sorted by their timestamp.
	LookupInstanceHistory(ctx context.Context, job Instance) (InstanceHistory, error)
	// ListInstanceHistory lists all state records for all job instances. The returned history is sorted by their timestamp.
	ListInstanceHistory(ctx context.Context) (InstanceHistory, error)
	// ClearInstanceHistory clears all state records for a job instance.
	ClearInstanceHistory(ctx context.Context) error
}

// LogStore is an interface that defines methods for logging job instance messages.
type LogStore interface {
	// Infof logs an informational message for a job instance.
	Infof(ctx context.Context, job Instance, format string, args ...any) error
	// Warnf logs a warning message for a job instance.
	Warnf(ctx context.Context, job Instance, format string, args ...any) error
	// Errorf logs an error message for a job instance.
	Errorf(ctx context.Context, job Instance, format string, args ...any) error
	// LookupInstanceLogs lists all log entries for a job instance. The returned logs are sorted by their timestamp.
	LookupInstanceLogs(ctx context.Context, job Instance) ([]Log, error)
	// ClearInstanceLogs clears all log entries for a job instance that match the specified filter.
	ClearInstanceLogs(ctx context.Context, filter Filter) error
}
