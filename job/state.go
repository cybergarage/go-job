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

	v1 "github.com/cybergarage/go-job/job/api/gen/go/v1"
)

// JobState represents the state of a job as an integer.
type JobState int

const (
	// JobStateUnset indicates that the job state is not set.
	JobStateUnset JobState = 0 // 0 indicates an unset state
	// JobCreated indicates the job has been created but not yet started.
	JobCreated JobState = 1 << iota
	// JobScheduled indicates the job has been scheduled for execution.
	JobScheduled
	// JobProcessing indicates the job is currently being processed.
	JobProcessing
	// JobCancelled indicates the job was cancelled before completion.
	JobCancelled
	// JobTimedOut indicates the job has exceeded its allowed execution time.
	JobTimedOut
	// JobCompleted indicates the job has completed (either successfully or unsuccessfully).
	JobCompleted
	// JobTerminated indicates the job has been terminated.
	JobTerminated
)

const (
	// JobStateActive represents the active states of a job (scheduled or processing).
	JobStateActive = JobScheduled | JobProcessing
	// JobStateFinal represents the final states of a job (cancelled, timed out, completed, or terminated).
	JobStateFinal = JobCancelled | JobTimedOut | JobCompleted | JobTerminated
	// JobStateError represents the error states of a job (cancelled, timed out, or terminated).
	JobStateError = JobCancelled | JobTimedOut | JobTerminated
	// JobStateSuccess represents the successful completion of a job.
	JobStateSuccess = JobCompleted
	// JobStateAll represents all possible states of a job.
	JobStateAll = JobCreated | JobScheduled | JobProcessing | JobCancelled | JobTimedOut | JobCompleted | JobTerminated
)

const (
	jobStateUnsetString      = "Unset"
	jobStateCreatedString    = "Created"
	jobStateScheduledString  = "Scheduled"
	jobStateProcessingString = "Processing"
	jobStateCancelledString  = "Cancelled"
	jobStateTimedOutString   = "TimedOut"
	jobStateCompletedString  = "Completed"
	jobStateTerminatedString = "Terminated"
)

// newStateFrom creates a new JobState from a given value.
func newStateFrom(a any) (JobState, error) {
	switch v := a.(type) {
	case JobState:
		return v, nil
	case string:
		return newStateFromString(v)
	case int:
		return JobState(v), nil
	case v1.JobState:
		switch v {
		case v1.JobState_JOB_STATE_UNSET:
			return JobStateUnset, nil
		case v1.JobState_JOB_STATE_CREATED:
			return JobCreated, nil
		case v1.JobState_JOB_STATE_SCHEDULED:
			return JobScheduled, nil
		case v1.JobState_JOB_STATE_PROCESSING:
			return JobProcessing, nil
		case v1.JobState_JOB_STATE_CANCELLED:
			return JobCancelled, nil
		case v1.JobState_JOB_STATE_TIMED_OUT:
			return JobTimedOut, nil
		case v1.JobState_JOB_STATE_COMPLETED:
			return JobCompleted, nil
		case v1.JobState_JOB_STATE_TERMINATED:
			return JobTerminated, nil
		}
	}
	return JobStateUnset, fmt.Errorf("invalid job state value: %v", a)
}

// newStateFromString returns the JobState corresponding to the given string.
func newStateFromString(s string) (JobState, error) {
	switch s {
	case jobStateCreatedString:
		return JobCreated, nil
	case jobStateScheduledString:
		return JobScheduled, nil
	case jobStateProcessingString:
		return JobProcessing, nil
	case jobStateCancelledString:
		return JobCancelled, nil
	case jobStateTimedOutString:
		return JobTimedOut, nil
	case jobStateCompletedString:
		return JobCompleted, nil
	case jobStateTerminatedString:
		return JobTerminated, nil
	case jobStateUnsetString:
		return JobStateUnset, nil
	default:
		return JobStateUnset, fmt.Errorf("unknown job state: %s", s)
	}
}

// Is checks if the current JobState is equal to the provided state.
func (s JobState) Is(state JobState) bool {
	return (s & state) != 0
}

// String returns the string representation of the JobState.
func (s JobState) String() string {
	switch s {
	case JobCreated:
		return jobStateCreatedString
	case JobScheduled:
		return jobStateScheduledString
	case JobProcessing:
		return jobStateProcessingString
	case JobCancelled:
		return jobStateCancelledString
	case JobTimedOut:
		return jobStateTimedOutString
	case JobCompleted:
		return jobStateCompletedString
	case JobTerminated:
		return jobStateTerminatedString
	default:
		return jobStateUnsetString
	}
}

// protoState converts the JobState to its corresponding protobuf representation.
func (s JobState) protoState() (v1.JobState, error) {
	switch s {
	case JobStateUnset:
		return v1.JobState_JOB_STATE_UNSET, nil
	case JobCreated:
		return v1.JobState_JOB_STATE_CREATED, nil
	case JobScheduled:
		return v1.JobState_JOB_STATE_SCHEDULED, nil
	case JobProcessing:
		return v1.JobState_JOB_STATE_PROCESSING, nil
	case JobCancelled:
		return v1.JobState_JOB_STATE_CANCELLED, nil
	case JobTimedOut:
		return v1.JobState_JOB_STATE_TIMED_OUT, nil
	case JobCompleted:
		return v1.JobState_JOB_STATE_COMPLETED, nil
	case JobTerminated:
		return v1.JobState_JOB_STATE_TERMINATED, nil
	}
	return v1.JobState_JOB_STATE_UNSET, fmt.Errorf("unknown job state: %s", s)
}
