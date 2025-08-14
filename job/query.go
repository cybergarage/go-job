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
	"github.com/google/uuid"
)

// Query is an interface that defines methods for querying job instances.
type Query interface {
	// UUID returns the UUID of the job instance.
	UUID() (uuid.UUID, bool)
	// Kind returns the kind of the job instance.
	Kind() (string, bool)
	// State returns the state of the job instance.
	State() (JobState, bool)
	// LogLevel returns the log level of the job instance.
	LogLevel() (LogLevel, bool)
	// IsAll returns true if the query matches all objects (no query criteria set).
	IsAll() bool
	// Matches checks if the specified object matches the query criteria.
	Matches(v any) bool
}

// QueryOption is a function that configures a job query.
type QueryOption func(*query)

type query struct {
	uuid  uuid.UUID
	kind  string
	state JobState
	level LogLevel
}

// WithQueryUUID sets the UUID for the job query.
func WithQueryUUID(uuid uuid.UUID) QueryOption {
	return func(q *query) {
		q.uuid = uuid
	}
}

// WithQueryKind sets the kind for the job query.
func WithQueryKind(kind string) QueryOption {
	return func(q *query) {
		q.kind = kind
	}
}

// WithQueryState sets the state for the job query.
func WithQueryState(state JobState) QueryOption {
	return func(q *query) {
		q.state = state
	}
}

// WithQueryLogLevel sets the level for the job query.
func WithQueryLogLevel(level LogLevel) QueryOption {
	return func(q *query) {
		q.level = level
	}
}

// WithQueryInstance sets the job query UUID and kind based on an existing job instance.
func WithQueryInstance(instance Instance) QueryOption {
	return func(q *query) {
		if instance == nil {
			return
		}
		q.uuid = instance.UUID()
		q.kind = instance.Kind()
	}
}

// NewQuery creates a new instance of Query with the given options.
func NewQuery(opts ...QueryOption) Query {
	q := &query{
		uuid:  uuid.Nil,
		kind:  "",
		state: JobStateUnset,
		level: LogNone,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// UUID returns the UUID of the job instance.
func (q *query) UUID() (uuid.UUID, bool) {
	if q.uuid == uuid.Nil {
		return uuid.Nil, false
	}
	return q.uuid, true
}

// Kind returns the kind of the job instance.
func (q *query) Kind() (string, bool) {
	if q.kind == "" {
		return "", false
	}
	return q.kind, true
}

// State returns the state of the job instance.
func (q *query) State() (JobState, bool) {
	if q.state == JobStateUnset {
		return JobStateUnset, false
	}
	return q.state, true
}

// LogLevel returns the log level of the job instance.
func (q *query) LogLevel() (LogLevel, bool) {
	if q.level == LogNone {
		return LogNone, false
	}
	return q.level, true
}

// IsAll returns true if the query matches all objects (no query criteria set).
func (q *query) IsAll() bool {
	if q == nil {
		return true
	}
	_, hasUUID := q.UUID()
	_, hasKind := q.Kind()
	_, hasState := q.State()
	_, hasLevel := q.LogLevel()
	return !hasUUID && !hasKind && !hasState && !hasLevel
}

// Matches checks if the specified object matches the query criteria.
func (q *query) Matches(v any) bool {
	if q.IsAll() {
		return true
	}
	switch v := v.(type) {
	case Instance:
		if uuid, ok := q.UUID(); ok && uuid != v.UUID() {
			return false
		}
		if kind, ok := q.Kind(); ok && kind != v.Kind() {
			return false
		}
		if state, ok := q.State(); ok && state != v.State() {
			return false
		}
		return true
	case InstanceState:
		if uuid, ok := q.UUID(); ok && uuid != v.UUID() {
			return false
		}
		if kind, ok := q.Kind(); ok && kind != v.Kind() {
			return false
		}
		if state, ok := q.State(); ok && !v.State().Matches(state) {
			return false
		}
		return true
	case Log:
		if uuid, ok := q.UUID(); ok && uuid != v.UUID() {
			return false
		}
		if kind, ok := q.Kind(); ok && kind != v.Kind() {
			return false
		}
		if level, ok := q.LogLevel(); ok && !v.Level().Contains(level) {
			return false
		}
		return true
	}
	return false
}
