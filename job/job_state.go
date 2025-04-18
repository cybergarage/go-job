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
	// JobCreated indicates the job has been created but not yet started.
	JobCreated JobState = iota

	// JobPending indicates the job is created but not yet started.
	JobPending

	// JobRunning indicates the job is currently executing.
	JobRunning // 1

	// JobSucceeded indicates the job has completed successfully.
	JobSucceeded // 2

	// JobFailed indicates the job has completed with an error.
	JobFailed // 3

	// JobCancelled indicates the job was cancelled before completion.
	JobCancelled // 4

	// JobTimedOut indicates the job has exceeded its allowed execution time.
	JobTimedOut // 5
)

// Is checks if the current JobState is equal to the provided state.
func (s JobState) Is(state JobState) bool {
	return s == state
}

// String returns the string representation of the JobState.
func (s JobState) String() string {
	switch s {
	case JobPending:
		return "Pending"
	case JobRunning:
		return "Running"
	case JobSucceeded:
		return "Succeeded"
	case JobFailed:
		return "Failed"
	case JobCancelled:
		return "Cancelled"
	case JobTimedOut:
		return "TimedOut"
	default:
		return "Unknown"
	}
}
