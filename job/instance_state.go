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
	"fmt"
	"time"

	"github.com/cybergarage/go-job/job/encoding"

	"github.com/google/uuid"
)

// InstanceState represents the state of a job instance at a specific point in time.
type InstanceState interface {
	// Kind returns the kind of the job instance.
	Kind() string
	// UUID returns the unique identifier of the job instance.
	UUID() uuid.UUID
	// Timestamp returns the timestamp of when the state history was created.
	Timestamp() time.Time
	// State returns the state of the job instance.
	State() JobState
	// Options returns the additional options associated with the instance record.
	Options() map[string]any
	// Map returns a map representation of the instance record.
	Map() map[string]any
	// JSONString returns a JSON string representation of the instance record.
	JSONString() (string, error)
	// String returns a string representation of the instance record.
	String() string
}

// InstanceStateOption is a function that configures the job instance state.
type InstanceStateOption func(*instanceState)

// WithStateKind is a functional option to set the kind of the instance state.
func WithStateKind(kind string) InstanceStateOption {
	return func(state *instanceState) {
		state.kind = kind
	}
}

// WithStateUUID is a functional option to set the UUID of the instance state.
func WithStateUUID(uuid uuid.UUID) InstanceStateOption {
	return func(state *instanceState) {
		state.uuid = uuid
	}
}

// WithStateJobState is a functional option to set the state of the instance state.
func WithStateJobState(s JobState) InstanceStateOption {
	return func(state *instanceState) {
		state.state = s
	}
}

// WithStateTimestamp is a functional option to set the timestamp of the instance state.
func WithStateTimestamp(ts time.Time) InstanceStateOption {
	return func(state *instanceState) {
		state.ts = ts
	}
}

// WithStateOption is a functional option to set additional options for the instance state.
func WithStateOption(opts map[string]any) func(*instanceState) {
	return func(state *instanceState) {
		for k, v := range opts {
			state.opts[k] = v
		}
	}
}

type instanceState struct {
	kind  string
	uuid  uuid.UUID
	ts    time.Time
	state JobState
	opts  map[string]any
}

// NewInstanceState creates a new job instance state with the provided options.
func NewInstanceState(opts ...InstanceStateOption) InstanceState {
	return newInstanceState(opts...)
}

// newInstanceState creates a new job state record with the current timestamp and the given state.
func newInstanceState(opts ...InstanceStateOption) InstanceState {
	state := &instanceState{
		kind:  "",
		uuid:  uuid.Nil,
		ts:    time.Now(),
		state: JobStateUnset,
		opts:  make(map[string]any),
	}
	for _, opt := range opts {
		opt(state)
	}
	return state
}

// NewInstanceStateFromMap creates a new instance state from a map representation.
func NewInstanceStateFromMap(m map[string]any) (InstanceState, error) {
	opts := make([]InstanceStateOption, 0, len(m))
	for key, value := range m {
		switch key {
		case kindKey:
			kind, err := NewKindFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithStateKind(kind))
		case uuidKey:
			uuid, err := NewUUIDFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithStateUUID(uuid))
		case timestampKey:
			ts, err := NewTimeFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithStateTimestamp(ts.Time()))
		case stateKey:
			state, err := NewStateFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithStateJobState(state))
		default:
			opts = append(opts, WithStateOption(map[string]any{key: value}))
		}
	}
	return newInstanceState(opts...), nil
}

// Kind returns the kind of the job instance.
func (state *instanceState) Kind() string {
	return state.kind
}

// UUID returns the unique identifier of the job instance.
func (state *instanceState) UUID() uuid.UUID {
	return state.uuid
}

// Timestamp returns the timestamp of when the state history was created.
func (state *instanceState) Timestamp() time.Time {
	return state.ts
}

// State returns the state of the job history.
func (state *instanceState) State() JobState {
	return state.state
}

// Options returns the additional options associated with the instance state.
func (state *instanceState) Options() map[string]any {
	return state.opts
}

// Map returns a map representation of the instance state.
func (state *instanceState) Map() map[string]any {
	m := map[string]any{
		kindKey:      state.kind,
		uuidKey:      state.uuid.String(),
		timestampKey: NewTimeFromTime(state.ts).String(),
		stateKey:     state.state.String(),
	}
	m = encoding.MergeMaps(m, state.opts)
	return m
}

// JSONString returns a JSON string representation of the instance state.
func (state *instanceState) JSONString() (string, error) {
	data, err := encoding.MapToJSON(state.Map())
	if err != nil {
		return "", err
	}
	return data, nil
}

// String returns a string representation of the instance state.
func (state *instanceState) String() string {
	return fmt.Sprintf("%s", state.Map())
}
