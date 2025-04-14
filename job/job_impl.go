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

type job struct {
	name        string
	uuid        uuid.UUID
	state       JobState
	payload     any
	handler     JobHandler
	createdAt   time.Time
	scheduledAt *time.Time
	startedAt   *time.Time
	finishedAt  *time.Time
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

// WithJobScheduledAt sets the scheduled time of the job.
func WithJobHandler(handler JobHandler) JobOption {
	return func(j *job) {
		j.handler = handler
	}
}

// WithJobScheduledAt sets the scheduled time of the job.
func WithJobScheduledAt(scheduledAt time.Time) JobOption {
	return func(j *job) {
		j.scheduledAt = &scheduledAt
	}
}

// NewJob creates a new job with the given name and options.
func NewJob(name string, opts ...JobOption) Job {
	j := &job{
		name:        name,
		uuid:        uuid.Nil,
		state:       JobCreated,
		payload:     nil,
		handler:     nil,
		createdAt:   time.Now(),
		scheduledAt: nil,
		startedAt:   nil,
		finishedAt:  nil,
	}

	for _, opt := range opts {
		opt(j)
	}

	if j.scheduledAt == nil {
		j.scheduledAt = &j.createdAt
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

// Handler returns the handler of the job.
func (j *job) Handler() JobHandler {
	return j.handler
}

// Payload returns the payload of the job.
func (j *job) Payload() any {
	return j.payload
}

// State returns the current state of the job.
func (j *job) State() JobState {
	return j.state
}

// SetState sets the state of the job.
func (j *job) SetState(state JobState) error {
	j.state = state
	return nil
}

// CreatedAt returns the creation time of the job.
func (j *job) CreatedAt() time.Time {
	return j.createdAt
}

// ScheduledAt returns the scheduled time of the job.
func (j *job) ScheduledAt() time.Time {
	return *j.scheduledAt
}

// StartedAt returns the start time of the job.
func (j *job) StartedAt() time.Time {
	return *j.startedAt
}

// FinishedAt returns the finish time of the job.
func (j *job) FinishedAt() time.Time {
	return *j.finishedAt
}

// String returns a string representation of the job.
func (j *job) String() string {
	return j.name
}
