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

// Log represents a log entry associated with a job.
type Log interface {
	// Kind returns the type of the log entry.
	Kind() string
	// UUID returns the unique identifier of the log entry.
	UUID() uuid.UUID
	// Timestamp returns the timestamp of the log entry.
	Timestamp() time.Time
	// Level returns the log level of the log entry.
	Level() LogLevel
	// Message returns the message of the log entry.
	Message() string
}

type log struct {
	kind  string
	uuid  uuid.UUID
	ts    time.Time
	level LogLevel
	msg   string
}

// NewLog creates a new log entry.
type LogOption func(*log)

// WithLogKind sets the type of the log entry.
func WithLogKind(kind string) LogOption {
	return func(l *log) {
		l.kind = kind
	}
}

// WithLogUUID sets the unique identifier of the log entry.
func WithLogUUID(uuid uuid.UUID) LogOption {
	return func(l *log) {
		l.uuid = uuid
	}
}

// WithLogLevel sets the log level of the log entry.
func WithLogLevel(level LogLevel) LogOption {
	return func(l *log) {
		l.level = level
	}
}

// WithLogMessage sets the message of the log entry.
func WithLogMessage(msg string) LogOption {
	return func(l *log) {
		l.msg = msg
	}
}

// WithLogTimestamp sets the timestamp of the log entry.
func WithLogTimestamp(ts time.Time) LogOption {
	return func(l *log) {
		l.ts = ts
	}
}

// NewLog creates a new log entry with the specified options.
func NewLog(opts ...LogOption) Log {
	l := &log{
		uuid:  uuid.Nil,
		ts:    time.Now(),
		level: LogInfo,
		msg:   "",
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Kind returns the type of the log entry.
func (l *log) Kind() string {
	return l.kind
}

// UUID returns the unique identifier of the log entry.
func (l *log) UUID() uuid.UUID {
	return l.uuid
}

// Timestamp returns the timestamp of the log entry.
func (l *log) Timestamp() time.Time {
	return l.ts
}

// Level returns the log level of the log entry.
func (l *log) Level() LogLevel {
	return l.level
}

// Message returns the message of the log entry.
func (l *log) Message() string {
	return l.msg
}

// String returns the string representation of the log entry.
func (l *log) String() string {
	return l.ts.Format(time.RFC3339) + " [" + l.level.String() + "] " + l.msg
}
