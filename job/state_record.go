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
)

// StateRecord keeps track of the state changes of a job.
type StateRecord struct {
	ts    time.Time
	state JobState
}

// NewStateRecord creates a new job state record with the current timestamp and the given state.
func NewStateRecord(state JobState) *StateRecord {
	return &StateRecord{
		ts:    time.Now(),
		state: state,
	}
}

// Timestamp returns the timestamp of when the state history was created.
func (sh *StateRecord) Timestamp() time.Time {
	return sh.ts
}

// State returns the state of the job history.
func (sh *StateRecord) State() JobState {
	return sh.state
}
