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

// InstanceRecord keeps track of the state changes of a job.
type InstanceRecord interface {
	UUID() uuid.UUID
	Timestamp() time.Time
	State() JobState
}

type instanceRecordOption func(*instanceRecord)

// WithRecordOption is a functional option to set additional options for the instance record.
func WithRecordOption(opts map[string]any) func(*instanceRecord) {
	return func(record *instanceRecord) {
		for k, v := range opts {
			record.opts[k] = v
		}
	}
}

type instanceRecord struct {
	id    uuid.UUID
	ts    time.Time
	state JobState
	opts  map[string]any
}

// newInstanceRecord creates a new job state record with the current timestamp and the given state.
func newInstanceRecord(id uuid.UUID, state JobState, opts ...instanceRecordOption) InstanceRecord {
	ir := &instanceRecord{
		id:    id,
		ts:    time.Now(),
		state: state,
		opts:  make(map[string]any),
	}
	for _, opt := range opts {
		opt(ir)
	}
	return ir
}

// UUID returns the unique identifier of the job instance.
func (record *instanceRecord) UUID() uuid.UUID {
	return record.id
}

// Timestamp returns the timestamp of when the state history was created.
func (record *instanceRecord) Timestamp() time.Time {
	return record.ts
}

// State returns the state of the job history.
func (record *instanceRecord) State() JobState {
	return record.state
}
