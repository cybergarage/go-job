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

// instanceRecord keeps track of the state changes of a job.
type instanceRecord struct {
	id    uuid.UUID
	ts    time.Time
	state JobState
}

// newInstanceRecord creates a new job state record with the current timestamp and the given state.
func newInstanceRecord(id uuid.UUID, state JobState) InstanceRecord {
	return &instanceRecord{
		id:    id,
		ts:    time.Now(),
		state: state,
	}
}

// UUID returns the unique identifier of the job instance.
func (sh *instanceRecord) UUID() uuid.UUID {
	return sh.id
}

// Timestamp returns the timestamp of when the state history was created.
func (sh *instanceRecord) Timestamp() time.Time {
	return sh.ts
}

// State returns the state of the job history.
func (sh *instanceRecord) State() JobState {
	return sh.state
}
