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
)

// Job represents a job that can be scheduled to run at a specific time or interval.
type Job interface {
	// Kind returns the name of the job.
	Kind() string
	// Description returns a description of the job.
	Description() string
	// Handler returns the job handler for the job.
	Handler() Handler
	// Schedule returns the schedule for the job.
	Schedule() Schedule
	// RegisteredAt returns the time when the job was registered.
	RegisteredAt() time.Time
	// Map returns a map representation of the job.
	Map() map[string]any
	// String returns a string representation of the job.
	String() string
}

type job struct {
	kind         string
	desc         string
	schedule     *schedule
	handler      *handler
	registeredAt time.Time
}

// JobOption is a function that configures a job.
type JobOption func(*job)

// WithKind sets the name of the job.
func WithKind(name string) JobOption {
	return func(j *job) {
		j.kind = name
	}
}

// WithDescription sets the description of the job.
func WithDescription(desc string) JobOption {
	return func(j *job) {
		j.desc = desc
	}
}

// NewJob creates a new job with the given name and options.
func NewJob(opts ...any) (Job, error) {
	return newJob(opts...)
}

func newJob(opts ...any) (*job, error) {
	schedule, err := newSchedule()
	if err != nil {
		return nil, err
	}

	j := &job{
		kind:         "",
		desc:         "",
		handler:      newHandler(),
		schedule:     schedule,
		registeredAt: time.Now(),
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case JobOption:
			opt(j)
		case HandlerOption:
			opt(j.handler)
		case ScheduleOption:
			opt(j.schedule)
		default:
			return nil, fmt.Errorf("invalid job option type: %T", opt)
		}
	}

	return j, nil
}

// Kind returns the name of the job.
func (j *job) Kind() string {
	return j.kind
}

// Description returns a description of the job.
func (j *job) Description() string {
	return j.desc
}

// Handler returns the handler of the job.
func (j *job) Handler() Handler {
	return j.handler
}

// RegisteredAt returns the time when the job was registered.
func (j *job) RegisteredAt() time.Time {
	return j.registeredAt
}

// Schedule returns the schedule of the job.
func (j *job) Schedule() Schedule {
	return j.schedule
}

// Map returns a map representation of the job.
func (j *job) Map() map[string]any {
	return map[string]any{
		kindKey: j.kind,
	}
}

// String returns a string representation of the job.
func (j *job) String() string {
	return fmt.Sprintf("%v", j.Map())
}
