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

	"github.com/cybergarage/go-job/job/encoding"
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
	// ScheduledAt returns the time when the job instance was scheduled.
	ScheduledAt() time.Time
	// Arguments returns the arguments for the job instance.
	Arguments() []any
	// Policy returns the policy associated with the job instance.
	Policy() Policy
	// UpdateState updates the state of the job instance and records the state change.
	UpdateState(state JobState, opts ...any) error
	// Process executes the job instance executor with the arguments provided in the context.
	Process() ([]any, error)
	// History returns the history of state changes for the job instance.
	History() (InstanceHistory, error)
	// State returns the current state of the job instance.
	State() JobState
	// AttemptCount returns the number of attempts made to process this job instance.
	AttemptCount() int
	// IsRecurring checks if the job instance is recurring.
	IsRecurring() bool
	// IsRetriable checks if the job instance can be retried.
	IsRetriable() bool
	// Map returns a map representation of the job instance.
	Map() map[string]any
	// String returns a string representation of the job instance.
	String() string
	// InstanceHelper provides methods to check if the job instance should be processed before or after another instance.
	InstanceHelper
}

type jobInstance struct {
	job     Job
	uuid    uuid.UUID
	state   JobState
	attempt int
	history History
	*handler
	*schedule
	*policy
	*arguments
}

// InstanceOption defines a function that configures a job instance.
type InstanceOption func(*jobInstance) error

// WithInstanceStore sets the history for the job instance.
func WithInstanceStore(history History) InstanceOption {
	return func(ji *jobInstance) error {
		ji.history = history
		return nil
	}
}

// NewInstance creates a new JobInstance with a unique identifier and initial state.
func NewInstance(opts ...any) (Instance, error) {
	ji := &jobInstance{
		job:       nil, // Default to nil, must be set by options
		uuid:      uuid.New(),
		state:     JobCreated,
		attempt:   0,
		history:   newHistory(),
		handler:   newHandler(),
		schedule:  newSchedule(),
		policy:    newPolicy(),
		arguments: newArguments(),
	}
	for _, opt := range opts {
		switch opt := opt.(type) {
		case InstanceOption:
			if err := opt(ji); err != nil {
				return nil, err
			}
		case HandlerOption:
			opt(ji.handler)
		case ScheduleOption:
			opt(ji.schedule)
		case PolicyOption:
			opt(ji.policy)
		case ArgumentsOption:
			opt(ji.arguments)
		case *arguments:
			ji.arguments = opt
		default:
			return nil, fmt.Errorf("invalid job instance option type: %T", opt)
		}
	}
	return ji, ji.UpdateState(JobCreated)
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

// Arguments returns the arguments for the job instance.
func (ji *jobInstance) Arguments() []any {
	if ji.arguments == nil {
		return nil
	}
	return ji.arguments.Arguments()
}

// Policy returns the policy associated with the job instance.
func (ji *jobInstance) Policy() Policy {
	return ji.policy
}

// ScheduledAt returns the time when the job instance was scheduled.
func (ji *jobInstance) ScheduledAt() time.Time {
	return ji.schedule.Next()
}

// Process executes the job instance executor with the arguments provided in the context.
func (ji *jobInstance) Process() ([]any, error) {
	if ji.handler == nil {
		return nil, fmt.Errorf("no job handler set for job instance %s", ji.uuid)
	}
	ji.attempt++
	res, err := ji.handler.Execute(ji.Arguments()...)
	if err == nil {
		ji.handler.HandleResponse(ji, res)
		return res, nil
	}
	return nil, ji.handler.HandleError(ji, err)
}

// UpdateState updates the state of the job instance and records the state change.
func (ji *jobInstance) UpdateState(state JobState, opts ...any) error {
	ji.state = state
	optMap := ji.OptionMap()
	for _, opt := range opts {
		switch opt := opt.(type) {
		case error:
			optMap["error"] = opt.Error()
		case Result:
			optMap["result"] = NewResultWith(opt).String()
		case map[string]any:
			optMap = encoding.MergeMaps(optMap, opt)
		default:
			return fmt.Errorf("invalid option type for UpdateState: %T", opt)
		}
	}
	return ji.history.LogInstanceRecord(ji, state, WithStateOption(optMap))
}

// History returns the history of state changes for the job instance.
func (ji *jobInstance) History() (InstanceHistory, error) {
	return ji.history.InstanceHistory(ji)
}

// State returns the current state of the job instance.
func (ji *jobInstance) State() JobState {
	return ji.state
}

// AttemptCount returns the number of attempts made to process this job instance.
func (ji *jobInstance) AttemptCount() int {
	return ji.attempt
}

// IsRetriable checks if the job instance can be retried based on its policy.
func (ji *jobInstance) IsRetriable() bool {
	maxRetries := ji.policy.MaxRetries()
	return maxRetries > 0 && ji.attempt < maxRetries
}

// Map returns a map representation of the job instance.
func (ji *jobInstance) Map() map[string]any {
	return encoding.MergeMaps(
		map[string]any{
			"uuid":  ji.uuid.String(),
			"state": ji.State().String(),
		},
		ji.OptionMap())
}

// OptionMap returns a map of options for the job instance, merging job, arguments, schedule, and policy options.
func (ji *jobInstance) OptionMap() map[string]any {
	mergedMap := map[string]any{}
	maps := []map[string]any{
		ji.arguments.Map(),
		ji.schedule.Map(),
		ji.Policy().Map(),
	}
	if ji.job != nil {
		maps = append(maps, ji.job.Map())
	}
	for _, m := range maps {
		mergedMap = encoding.MergeMaps(mergedMap, m)
	}
	return mergedMap
}

// String returns a string representation of the job instance.
func (ji *jobInstance) String() string {
	return fmt.Sprintf("%v", ji.Map())
}
