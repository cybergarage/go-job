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

// JobHistory keeps track of the state changes of a job.
type JobHistory struct {
	ts    time.Time
	state JobState
}

func NewJobHistory(state JobState) *JobHistory {
	return &JobHistory{
		ts:    time.Now(),
		state: state,
	}
}

// Timestamp returns the timestamp of when the job history was created.
func (jh *JobHistory) Timestamp() time.Time {
	return jh.ts
}

// State returns the state of the job history.
func (jh *JobHistory) State() JobState {
	return jh.state
}
