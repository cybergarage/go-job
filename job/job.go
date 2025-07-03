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
)

// Job represents a job that can be scheduled to run at a specific time or interval.
type Job interface {
	// Kind returns the name of the job.
	Kind() string
	// Handler returns the job handler for the job.
	Handler() JobHandler
	// String returns a string representation of the job.
	String() string
}

type job struct {
	name    string
	logger  Logger
	handler *jobHandler
	ctx     *ctx
}

// JobOption is a function that configures a job.
type JobOption func(*job)

// WithJobName sets the name of the job.
func WithJobName(name string) JobOption {
	return func(j *job) {
		j.name = name
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
		case ContextOption:
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

// String returns a string representation of the job.
func (j *job) String() string {
	return j.name
}
