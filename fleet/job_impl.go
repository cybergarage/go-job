// Copyright (C) 2025 The go-fleet Authors. All rights reserved.
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

package fleet

import (
	"github.com/google/uuid"
)

type job struct {
	name    string
	uuid    uuid.UUID
	state   JobState
	payload any
}

// JobOption is a function that configures a job.
type JobOption func(*job)

// WithJobName sets the name of the job.
func WithJobName(name string) JobOption {
	return func(j *job) {
		j.name = name
	}
}

// WithJobState sets the state of the job.
func WithJobState(state JobState) JobOption {
	return func(j *job) {
		j.state = state
	}
}

// WithJobPayload sets the payload of the job.
func WithJobPayload(payload any) JobOption {
	return func(j *job) {
		j.payload = payload
	}
}

// WithJobUUID sets the UUID of the job.
// If you want to generate a new UUID, use the NewJob function.
func WithJobUUID(uuid uuid.UUID) JobOption {
	return func(j *job) {
		j.uuid = uuid
	}
}

// NewJob creates a new job with the given name and options.
func NewJob(name string, opts ...JobOption) Job {
	j := &job{
		name:    name,
		uuid:    uuid.Nil,
		state:   Pending,
		payload: nil,
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

// Name returns the name of the job.
func (j *job) Name() string {
	return j.name
}

// UUID returns the UUID of the job. If the UUID is not set, it generates a new one.
func (j *job) UUID() uuid.UUID {
	if j.uuid == uuid.Nil {
		j.uuid = uuid.New()
	}
	return j.uuid
}

// Payload returns the payload of the job.
func (j *job) Payload() any {
	return j.payload
}

// State returns the current state of the job.
func (j *job) State() JobState {
	return j.state
}

// String returns a string representation of the job.
func (j *job) String() string {
	return j.name
}
