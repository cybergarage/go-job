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

// JobState represents the state of a job as an integer.
type JobState int

const (
	// JobStateUnknown indicates an unknown state, typically used for uninitialized jobs.
	JobStateUnknown JobState = 0 // 0 indicates an unknown state
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

// Is checks if the current JobState is equal to the provided state.
func (s JobState) Is(state JobState) bool {
	return s == state
}

// String returns the string representation of the JobState.
func (s JobState) String() string {
	switch s {
	case JobCreated:
		return "Created"
	case JobScheduled:
		return "Scheduled"
	case JobProcessing:
		return "Processing"
	case JobCancelled:
		return "Cancelled"
	case JobTimedOut:
		return "TimedOut"
	case JobCompleted:
		return "Completed"
	case JobTerminated:
		return "Terminated"
	default:
		return "Unknown"
	}
}
