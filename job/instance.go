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
	// Schedule defines the scheduling interface for the job instance.
	Schedule
	// Handler defines the handler interface for the job instance.
	Handler
	// Policy defines the policy interface for the job instance.
	Policy
	// Arguments defines the arguments interface for the job instance.
	Arguments

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
	// Handler returns the handler for the job instance.
	Handler() Handler
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
	// JSONString returns a JSON string representation of the job instance.
	JSONString() (string, error)
	// String returns a string representation of the job instance.
	String() string
	// InstanceLogger provides methods for logging messages related to the job instance.
	instanceLogger
	// InstanceHelper provides methods to check if the job instance should be processed before or after another instance.
	instanceHelper
}

type jobInstance struct {
	*handler
	*schedule
	*policy
	*argumentsImpl

	job          *job
	uuid         UUID
	state        JobState
	attempt      int
	history      History
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

// WithJob sets the job for the job instance.
func WithJob(job Job) InstanceOption {
	return func(ji *jobInstance) error {
		jobOpts := []JobOption{
			WithKind(job.Kind()),
			WithDescription(job.Description()),
			WithRegisteredAt(job.RegisteredAt()),
		}
		for _, opt := range jobOpts {
			opt(ji.job)
		}

		handlerOpts := []HandlerOption{
			WithExecutor(job.Handler().Executor()),
			WithStateChangeProcessor(job.Handler().StateChangeProcessor()),
			WithTerminateProcessor(job.Handler().TerminateProcessor()),
			WithCompleteProcessor(job.Handler().CompleteProcessor()),
		}
		for _, opt := range handlerOpts {
			opt(ji.handler)
		}

		scheduleOpts := []ScheduleOption{
			WithCrontabSpec(job.Schedule().CrontabSpec()),
			WithScheduleAt(job.Schedule().Next()),
		}
		for _, opt := range scheduleOpts {
			if err := opt(ji.schedule); err != nil {
				return err
			}
		}

		return nil
	}
}

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

// WithInstanceStore sets the store for the job instance history.
func WithInstanceStore(store Store) InstanceOption {
	return func(ji *jobInstance) error {
		ji.history = newHistory(
			withHistoryStore(store),
		)
		return nil
	}
}

// WithAttemptCount sets the number of attempts made to process the job instance.
func WithAttemptCount(attempt int) InstanceOption {
	return func(ji *jobInstance) error {
		ji.attempt = attempt
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

	schedule, err := newSchedule(WithScheduleAt(time.Now()))
	if err != nil {
		return nil, err
	}

	ji := &jobInstance{
		job:           job,
		uuid:          NewUUID(),
		state:         JobStateUnset,
		attempt:       0,
		history:       newHistory(),
		handler:       newHandler(),
		schedule:      schedule,
		policy:        newPolicy(),
		argumentsImpl: newArguments(),
		createdAt:     time.Now(),
		completedAt:   time.Time{},
		terminatedAt:  time.Time{},
		processedAt:   time.Time{},
		cancelledAt:   time.Time{},
		timedoutAt:    time.Time{},
		resultSet:     nil,
		resultError:   nil,
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
			opt(ji.argumentsImpl)
		case *argumentsImpl:
			ji.argumentsImpl = opt
		default:
			return nil, fmt.Errorf("invalid job instance option type: %T", opt)
		}
	}

	return ji, nil
}

// NewInstanceFromMap creates a new job instance from the provided map and options.
func NewInstanceFromMap(m map[string]any, opts ...any) (Instance, error) {
	for key, value := range m {
		switch key {
		case kindKey:
			kind, err := newKindFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithKind(kind))
		case uuidKey:
			uuid, err := NewUUIDFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithUUID(uuid))
		case stateKey:
			state, err := newStateFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithState(state))
		case argumentsKey:
			args, err := newArgumentsFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithArguments(args.Arguments()...))
		case crontabKey:
			crontabSpec, err := newCrontabSpecFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithCrontabSpec(crontabSpec))
		case scheduleAtKey:
			scheduleAt, err := NewTimestampFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithScheduleAt(scheduleAt.Time()))
		case maxRetriesKey:
			maxRetries, err := newMaxRetriesFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithMaxRetries(maxRetries))
		case priorityKey:
			priority, err := NewPriorityFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithPriority(priority))
		case timeoutKey:
			timeout, err := newTimeoutFrom(value)
			if err != nil {
				return nil, err
			}
			opts = append(opts, WithTimeout(timeout))
		}
	}
	return NewInstance(opts...)
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
	if ji.argumentsImpl == nil {
		return nil
	}
	return ji.argumentsImpl.Arguments()
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
	return ji.Next()
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

// Handler returns the handler for the job instance.
func (ji *jobInstance) Handler() Handler {
	return ji.handler
}

// Process executes the job instance executor with the arguments provided in the context.
func (ji *jobInstance) Process() ([]any, error) {
	if ji.handler == nil {
		return nil, fmt.Errorf("no job handler set for job instance %s", ji.uuid)
	}
	ji.attempt++
	ji.resultSet, ji.resultError = ji.Execute(ji.Arguments()...)
	if ji.resultError != nil {
		ji.resultError = ji.HandleTerminated(ji, ji.resultError)
	}
	return ji.resultSet, ji.resultError
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

	opts = append(opts,
		ji.argumentsImpl.Map(),
	)

	optMap := ji.OptionMap()
	for _, opt := range opts {
		switch opt := opt.(type) {
		case error:
			optMap[errorKey] = opt.Error()
		case ResultSet:
			optMap[resultSetKey] = newResultWith(opt).String()
		case map[string]any:
			optMap = encoding.MergeMaps(optMap, opt)
		}
	}

	err := ji.history.LogProcessState(ji, state, withStateOption(optMap))
	if err != nil {
		return err
	}

	chgProcessor := ji.Handler().StateChangeProcessor()
	if chgProcessor != nil {
		chgProcessor(ji, state)
	}

	return nil
}

// History returns the history of state changes for the job instance.
func (ji *jobInstance) History() (InstanceHistory, error) {
	return ji.history.LookupHistory(NewQuery(WithQueryInstance(ji)))
}

// Logs returns the logs for the job instance.
func (ji *jobInstance) Logs() ([]Log, error) {
	return ji.history.LookupLogs(NewQuery(WithQueryInstance(ji)))
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
	maxRetries := ji.MaxRetries()
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
		ji.argumentsImpl.Map(),
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
