// Copyright (C) 2025 The go-fleet Authors. All rights reserved.
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

package fleet

// JobState represents the state of a job as an integer.
type JobState int

const (
	// Pending indicates the job is created but not yet started.
	Pending JobState = iota // 0

	// Running indicates the job is currently executing.
	Running // 1

	// Succeeded indicates the job has completed successfully.
	Succeeded // 2

	// Failed indicates the job has completed with an error.
	Failed // 3

	// Cancelled indicates the job was cancelled before completion.
	Cancelled // 4

	// TimedOut indicates the job has exceeded its allowed execution time.
	TimedOut // 5
)

// Is checks if the current JobState is equal to the provided state.
func (s JobState) Is(state JobState) bool {
	return s == state
}

// String returns the string representation of the JobState.
func (s JobState) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Running:
		return "Running"
	case Succeeded:
		return "Succeeded"
	case Failed:
		return "Failed"
	case Cancelled:
		return "Cancelled"
	case TimedOut:
		return "TimedOut"
	default:
		return "Unknown"
	}
}
