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
	// CreatedAt returns the time when the job instance was created.
	CreatedAt() time.Time
	// ScheduledAt returns the time when the job instance was scheduled.
	ScheduledAt() time.Time
	// ProcessedAt returns the time when the job instance started processing.
	ProcessedAt() time.Time
	// CompletedAt returns the completed time when the job instance was completed.
	CompletedAt() time.Time
	// TerminatedAt returns the terminated time when the job instance was terminated.
	TerminatedAt() time.Time
	// CancelledAt returns the time when the job instance was cancelled.
	CancelledAt() time.Time
	// TimeoutedAt returns the time when the job instance timed out.
	TimeoutedAt() time.Time
	// Arguments returns the arguments for the job instance.
	Arguments() []any
	// Policy returns the policy associated with the job instance.
	Policy() Policy
	// UpdateState updates the state of the job instance and records the state change.
	UpdateState(state JobState, opts ...any) error
	// Process executes the job instance executor with the arguments provided in the context.
	Process() ([]any, error)
	// Result returns the processed result set of the executor when the job instance is completed or terminated.
	// If the job instance is not completed or terminated, it returns an error.
	ResultSet() (ResultSet, error)
	// History returns the history of state changes for the job instance.
	History() (InstanceHistory, error)
	// Logs returns the logs for the job instance.
	Logs() ([]Log, error)
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
	// InstanceLogger provides methods for logging messages related to the job instance.
	InstanceLogger
	// InstanceHelper provides methods to check if the job instance should be processed before or after another instance.
	InstanceHelper
}

type jobInstance struct {
	job     *job
	uuid    uuid.UUID
	state   JobState
	attempt int
	history History
	*handler
	*schedule
	*policy
	*arguments
	createdAt    time.Time
	completedAt  time.Time
	terminatedAt time.Time
	processedAt  time.Time
	cancelledAt  time.Time
	timedoutAt   time.Time
	resultSet    ResultSet
	resultError  error
}

// InstanceOption defines a function that configures a job instance.
type InstanceOption func(*jobInstance) error

// WithUUID sets the unique identifier for the job instance.
func WithUUID(uuid uuid.UUID) InstanceOption {
	return func(ji *jobInstance) error {
		ji.uuid = uuid
		return nil
	}
}

// WithInstanceHistory sets the history for the job instance.
func WithInstanceHistory(history History) InstanceOption {
	return func(ji *jobInstance) error {
		ji.history = history
		return nil
	}
}

// WithCreatedAt sets the time when the job instance was created.
func WithCreatedAt(t time.Time) InstanceOption {
	return func(ji *jobInstance) error {
		ji.createdAt = t
		return nil
	}
}

// WithProcessingAt sets the time when the job instance started processing.
func WithProcessingAt(t time.Time) InstanceOption {
	return func(ji *jobInstance) error {
		ji.processedAt = t
		return nil
	}
}

// WithCompletedAt sets the time when the job instance was completed.
func WithCompletedAt(t time.Time) InstanceOption {
	return func(ji *jobInstance) error {
		ji.completedAt = t
		return nil
	}
}

// WithTerminatedAt sets the time when the job instance was terminated.
func WithTerminatedAt(t time.Time) InstanceOption {
	return func(ji *jobInstance) error {
		ji.terminatedAt = t
		return nil
	}
}

// WithCancelledAt sets the time when the job instance was cancelled.
func WithCancelledAt(t time.Time) InstanceOption {
	return func(ji *jobInstance) error {
		ji.cancelledAt = t
		return nil
	}
}

// WithTimedOutAt sets the time when the job instance timed out.
func WithTimedOutAt(t time.Time) InstanceOption {
	return func(ji *jobInstance) error {
		ji.timedoutAt = t
		return nil
	}
}

// WithResultSet sets the result set for the job instance.
func WithResultSet(rs ResultSet) InstanceOption {
	return func(ji *jobInstance) error {
		ji.resultSet = rs
		return nil
	}
}

// WithResultError sets the error for the job instance result.
func WithResultError(err error) InstanceOption {
	return func(ji *jobInstance) error {
		ji.resultError = err
		return nil
	}
}

// WithState sets the state of the job instance.
func WithState(state JobState) InstanceOption {
	return func(ji *jobInstance) error {
		ji.state = state
		return nil
	}
}

// NewInstance creates a new JobInstance with a unique identifier and initial state.
func NewInstance(opts ...any) (Instance, error) {
	job, err := newJob()
	if err != nil {
		return nil, err
	}
	ji := &jobInstance{
		job:          job,
		uuid:         uuid.New(),
		state:        JobStateUnset,
		attempt:      0,
		history:      newHistory(),
		handler:      newHandler(),
		schedule:     newSchedule(),
		policy:       newPolicy(),
		arguments:    newArguments(),
		createdAt:    time.Now(),
		completedAt:  time.Time{},
		terminatedAt: time.Time{},
		processedAt:  time.Time{},
		cancelledAt:  time.Time{},
		timedoutAt:   time.Time{},
		resultSet:    nil,
		resultError:  nil,
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case InstanceOption:
			if err := opt(ji); err != nil {
				return nil, err
			}
		case JobOption:
			opt(ji.job)
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

// CreatedAt returns the time when the job instance was created.
func (ji *jobInstance) CreatedAt() time.Time {
	return ji.createdAt
}

// ScheduledAt returns the time when the job instance was scheduled.
func (ji *jobInstance) ScheduledAt() time.Time {
	return ji.schedule.Next()
}

// ProcessedAt returns the time when the job instance started processing.
func (ji *jobInstance) ProcessedAt() time.Time {
	return ji.processedAt
}

// CompletedAt returns the completed time when the job instance was completed.
func (ji *jobInstance) CompletedAt() time.Time {
	return ji.completedAt
}

// TerminatedAt returns the terminated time when the job instance was terminated.
func (ji *jobInstance) TerminatedAt() time.Time {
	return ji.terminatedAt
}

// CancelledAt returns the time when the job instance was cancelled.
func (ji *jobInstance) CancelledAt() time.Time {
	return ji.cancelledAt
}

// TimeoutedAt returns the time when the job instance timed out.
func (ji *jobInstance) TimeoutedAt() time.Time {
	return ji.timedoutAt
}

// Process executes the job instance executor with the arguments provided in the context.
func (ji *jobInstance) Process() ([]any, error) {
	if ji.handler == nil {
		return nil, fmt.Errorf("no job handler set for job instance %s", ji.uuid)
	}
	ji.attempt++
	ji.resultSet, ji.resultError = ji.handler.Execute(ji.Arguments()...)
	if ji.resultError == nil {
		ji.handler.HandleResponse(ji, ji.resultSet)
		return ji.resultSet, nil
	}
	ji.resultError = ji.handler.HandleError(ji, ji.resultError)
	return nil, ji.resultError
}

// Result returns the processed result set of the executor when the job instance is completed or terminated.
// If the job instance is not completed or terminated, it returns an error.
func (ji *jobInstance) ResultSet() (ResultSet, error) {
	if ji.state != JobCompleted && ji.state != JobTerminated {
		return nil, fmt.Errorf("job instance %s is not completed or terminated, current state: %s", ji.uuid, ji.state)
	}
	return ji.resultSet, ji.resultError
}

// UpdateState updates the state of the job instance and records the state change.
func (ji *jobInstance) UpdateState(state JobState, opts ...any) error {
	ji.state = state
	switch state {
	case JobCompleted:
		ji.completedAt = time.Now()
	case JobTerminated:
		ji.terminatedAt = time.Now()
	}

	optMap := ji.OptionMap()
	for _, opt := range opts {
		switch opt := opt.(type) {
		case error:
			optMap[errorKey] = opt.Error()
		case ResultSet:
			optMap[resultSetKey] = NewResultWith(opt).String()
		case map[string]any:
			optMap = encoding.MergeMaps(optMap, opt)
		}
	}

	err := ji.history.LogProcessState(ji, state, WithStateOption(optMap))
	if err != nil {
		return err
	}

	return nil
}

// History returns the history of state changes for the job instance.
func (ji *jobInstance) History() (InstanceHistory, error) {
	return ji.history.LookupHistory(ji)
}

// Logs returns the logs for the job instance.
func (ji *jobInstance) Logs() ([]Log, error) {
	return ji.history.LookupLogs(ji)
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
			kindKey:  ji.Kind(),
			uuidKey:  ji.uuid.String(),
			stateKey: ji.State().String(),
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
