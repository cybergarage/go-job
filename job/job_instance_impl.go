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

	"github.com/google/uuid"
)

type jobInstance struct {
	job  Job
	uuid uuid.UUID
}

// JobInstanceOption defines a function that configures a job instance.
type JobInstanceOption func(*jobInstance) error

// WithJobInstanceJob sets the job for the job instance.
func WithJobInstanceJob(job Job) JobInstanceOption {
	return func(ji *jobInstance) error {
		ji.job = job
		return nil
	}
}

// NewJobInstance creates a new JobInstance with a unique identifier and initial state.
func NewJobInstance(opts ...JobInstanceOption) (JobInstance, error) {
	ji := &jobInstance{
		uuid: uuid.New(),
	}
	for _, opt := range opts {
		if err := opt(ji); err != nil {
			return nil, err
		}
	}
	return ji, nil
}

// UUID returns the unique identifier of the job instance.
func (ji *jobInstance) UUID() uuid.UUID {
	return ji.uuid
}

// State returns the current state of the job instance.
func (ji *jobInstance) State() JobState {
	return JobPending
}

// String returns a string representation of the job instance.
func (ji *jobInstance) String() string {
	return fmt.Sprintf("JobInstance{UUID: %s, Job: %v, State: %v}", ji.uuid, ji.job, ji.State())
}
