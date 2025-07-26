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
	// Equal checks if two log entries are equal.
	Equal(other Log) bool
	// Map returns a map representation of the log entry.
	Map() map[string]any
	// String returns the string representation of the log entry.
	String() string
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

// NewLogFromMap creates a new log entry from a map representation.
func NewLogFromMap(m map[string]any) (Log, error) {
	opts := []LogOption{}
	for key, value := range m {
		switch key {
		case kindKey:
			kind, err := NewKindFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithLogKind(kind))
		case uuidKey:
			uuid, err := NewUUIDFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithLogUUID(uuid))
		case timestampKey:
			ts, err := NewTimestampFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithLogTimestamp(ts.Time()))
		case levelKey:
			level, err := NewLogLevelFromString(value.(string))
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithLogLevel(level))
		case messageKey:
			if msg, ok := value.(string); ok {
				opts = append(opts, WithLogMessage(msg))
			}
		}
	}
	return NewLog(opts...), nil
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

// Equal checks if two log entries are equal.
func (l *log) Equal(other Log) bool {
	return l.Kind() == other.Kind() &&
		l.UUID() == other.UUID() &&
		l.Timestamp().Equal(other.Timestamp()) &&
		l.Level() == other.Level() &&
		l.Message() == other.Message()
}

// Map returns a map representation of the log entry.
func (l *log) Map() map[string]any {
	return map[string]any{
		kindKey:      l.kind,
		uuidKey:      l.uuid.String(),
		timestampKey: NewTimestampFromTime(l.ts).String(),
		levelKey:     l.level.String(),
		messageKey:   l.msg,
	}
}

// String returns the string representation of the log entry.
func (l *log) String() string {
	return NewTimestampFromTime(l.ts).String() + " [" + l.level.String() + "] " + l.msg
}
