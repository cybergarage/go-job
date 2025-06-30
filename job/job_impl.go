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

import "fmt"

type job struct {
	name    string
	payload any
	logger  Logger
	handler *jobHandler
	ctx     *jobContext
}

// JobOption is a function that configures a job.
type JobOption func(*job)

// WithJobName sets the name of the job.
func WithJobName(name string) JobOption {
	return func(j *job) {
		j.name = name
	}
}

// WithJobPayload sets the payload of the job.
func WithJobPayload(payload any) JobOption {
	return func(j *job) {
		j.payload = payload
	}
}

// WithJobStartedAt sets the start time of the job.
func WithJobLogger(logger Logger) JobOption {
	return func(j *job) {
		j.logger = logger
	}
}

// NewJob creates a new job with the given name and options.
func NewJob(opts ...any) (Job, error) {
	j := &job{
		name:    "",
		payload: nil,
		handler: newJobHandler(),
		ctx:     newJobContext(),
		logger:  NewNullLogger(),
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case JobOption:
			opt(j)
		case JobHandlerOption:
			opt(j.handler)
		case JobContextOption:
			opt(j.ctx)
		default:
			return nil, fmt.Errorf("invalid job option type: %T", opt)
		}
	}

	return j, nil
}

// Kind returns the name of the job.
func (j *job) Kind() string {
	return j.name
}

// Handler returns the handler of the job.
func (j *job) Handler() JobHandler {
	return j.handler
}

// Payload returns the payload of the job.
func (j *job) Payload() any {
	return j.payload
}

func (j *job) Process() error {
	// Implement the logic to process the job
	// For example, you might want to execute the job's command or function
	// and update its state in the database or in-memory store.
	return nil
}

// String returns a string representation of the job.
func (j *job) String() string {
	return j.name
}
