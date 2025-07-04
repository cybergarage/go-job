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

// Instance represents a specific instance of a job that has been scheduled or executed.
type Instance interface {
	// Job returns the job associated with this job instance.
	Job() Job
	// Kind returns the kind of job this instance represents.
	Kind() Kind
	// UUID returns the unique identifier of the job instance.
	UUID() uuid.UUID
	// UpdateState updates the state of the job instance and records the state change.
	UpdateState(state JobState) error
	// Process executes the job instance executor with the arguments provided in the context.
	Process() error
	// State returns the current state of the job instance.
	State() JobState
	// String returns a string representation of the job instance.
	String() string
}
type jobInstance struct {
	job  Job
	uuid uuid.UUID
	*StateHistory
	handler *handler
	args    *Arguments
}

// InstanceOption defines a function that configures a job instance.
type InstanceOption func(*jobInstance) error

// NewInstance creates a new JobInstance with a unique identifier and initial state.
func NewInstance(opts ...any) (Instance, error) {
	ji := &jobInstance{
		job:          nil, // Default to nil, must be set by options
		uuid:         uuid.New(),
		StateHistory: NewStateHistory(),
		handler:      newHandler(),
		args:         NewArguments(),
	}
	for _, opt := range opts {
		switch opt := opt.(type) {
		case InstanceOption:
			if err := opt(ji); err != nil {
				return nil, err
			}
		case HandlerOption:
			opt(ji.handler)
		case ArgumentsOption:
			opt(ji.args)
		case *Arguments:
			ji.args = opt
		default:
			return nil, fmt.Errorf("invalid job instance option type: %T", opt)
		}
	}
	return ji, nil
}

// Job returns the job associated with this job instance.
func (ji *jobInstance) Job() Job {
	return ji.job
}

// Kind returns the kind of job this instance represents.
func (ji *jobInstance) Kind() Kind {
	if ji.job == nil {
		return ""
	}
	return ji.job.Kind()
}

// UUID returns the unique identifier of the job instance.
func (ji *jobInstance) UUID() uuid.UUID {
	return ji.uuid
}

// UpdateState updates the state of the job instance and records the state change.
func (ji *jobInstance) UpdateState(state JobState) error {
	ji.AppendStateRecord(state)
	return nil
}

// Process executes the job instance executor with the arguments provided in the context.
func (ji *jobInstance) Process() error {
	if ji.handler == nil {
		return fmt.Errorf("no job handler set for job instance %s", ji.uuid)
	}
	res, err := ji.handler.Execute(ji.args.Args()...)
	if err == nil {
		ji.handler.HandleResponse(ji, res)
		return nil
	}
	return ji.handler.HandleError(ji, err)
}

// State returns the current state of the job instance.
func (ji *jobInstance) State() JobState {
	r := ji.LastStateRecord()
	if r == nil {
		return JobStateUnknown
	}
	return r.State()
}

// String returns a string representation of the job instance.
func (ji *jobInstance) String() string {
	return fmt.Sprintf("JobInstance{UUID: %s, Job: %v, State: %v}", ji.uuid, ji.job, ji.State())
}
