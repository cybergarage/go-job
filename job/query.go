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
	"time"

	"github.com/google/uuid"
)

// Query defines an interface for specifying and evaluating query criteria for jobs, instances, and logs.
type Query interface {
	// Filter provides additional filtering methods for time-based or custom criteria.
	Filter

	// UUID returns the UUID criterion for the query, if set.
	UUID() (uuid.UUID, bool)
	// Kind returns the kind criterion for the query, if set.
	Kind() (string, bool)
	// State returns the job state criterion for the query, if set.
	State() (JobState, bool)
	// LogLevel returns the log level criterion for the query, if set.
	LogLevel() (LogLevel, bool)

	// IsUnset returns true if no query criteria are set.
	IsUnset() bool

	// Matches returns true if the specified object satisfies all query criteria.
	Matches(v any) bool
}

// QueryOption is a function that configures a job query.
type QueryOption func(*query)

type query struct {
	*filter

	uuid  uuid.UUID
	kind  string
	state JobState
	level LogLevel
}

// WithQueryUUID sets the UUID for the query.
func WithQueryUUID(uuid uuid.UUID) QueryOption {
	return func(q *query) {
		q.uuid = uuid
	}
}

// WithQueryKind sets the kind for the query.
func WithQueryKind(kind string) QueryOption {
	return func(q *query) {
		q.kind = kind
	}
}

// WithQueryState sets the state for the query.
func WithQueryState(state JobState) QueryOption {
	return func(q *query) {
		q.state = state
	}
}

// WithQueryLogLevel sets the level for the query.
func WithQueryLogLevel(level LogLevel) QueryOption {
	return func(q *query) {
		q.level = level
	}
}

// WithQueryInstance sets the query UUID and kind based on an existing job instance.
func WithQueryInstance(instance Instance) QueryOption {
	return func(q *query) {
		if instance == nil {
			return
		}
		q.uuid = instance.UUID()
		q.kind = instance.Kind()
	}
}

// WithQueryBefore sets the time before which target should be filtered.
func WithQueryBefore(before time.Time) QueryOption {
	return func(q *query) {
		q.filter.before = before
	}
}

// WithQueryAfter sets the time after which target should be filtered.
func WithQueryAfter(after time.Time) QueryOption {
	return func(q *query) {
		q.filter.after = after
	}
}

func NewQuery(opts ...QueryOption) Query {
	q := &query{
		filter: newFilter(),
		uuid:   uuid.Nil,
		kind:   "",
		state:  JobStateUnset,
		level:  LogNone,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// UUID returns the UUID criterion for the query, if set.
func (q *query) UUID() (uuid.UUID, bool) {
	if q.uuid == uuid.Nil {
		return uuid.Nil, false
	}
	return q.uuid, true
}

// Kind returns the kind criterion for the query, if set.
func (q *query) Kind() (string, bool) {
	if q.kind == "" {
		return "", false
	}
	return q.kind, true
}

// State returns the job state criterion for the query, if set.
func (q *query) State() (JobState, bool) {
	if q.state == JobStateUnset {
		return JobStateUnset, false
	}
	return q.state, true
}

// LogLevel returns the log level criterion for the query, if set.
func (q *query) LogLevel() (LogLevel, bool) {
	if q.level == LogNone {
		return LogNone, false
	}
	return q.level, true
}

// IsUnset returns true if no query criteria are set.
func (q *query) IsUnset() bool {
	if q == nil {
		return true
	}
	_, hasUUID := q.UUID()
	_, hasKind := q.Kind()
	_, hasState := q.State()
	_, hasLevel := q.LogLevel()
	hasFilter := !q.filter.IsUnset()
	return !hasUUID && !hasKind && !hasState && !hasLevel && !hasFilter
}

// Matches returns true if the specified object satisfies all query criteria.
func (q *query) Matches(v any) bool {
	if q.IsUnset() {
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
		if ok := !q.filter.IsUnset(); ok && !q.filter.Matches(v) {
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
		if ok := !q.filter.IsUnset(); ok && !q.filter.Matches(v) {
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
		if ok := !q.filter.IsUnset(); ok && !q.filter.Matches(v) {
			return false
		}
		return true
	}
	return false
}
