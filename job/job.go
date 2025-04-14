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
	"context"

	"github.com/google/uuid"
)

// JobHandler is an interface that defines methods for handling jobs.
type JobHandler interface {
	// ProcessJob processes a job by executing it and updating its state.
	ProcessJob(ctx context.Context, job Job) error
	// HandleJobError handles errors that occur during job processing.
	HandleJobError(ctx context.Context, job Job, err error) error
}

// Job represents a job that can be scheduled to run at a specific time or interval.
type Job interface {
	// Name returns the name of the job.
	Name() string
	// UUID returns the UUID of the job.
	UUID() uuid.UUID
	// Payload returns the payload of the job.
	Payload() any
	// State returns the current state of the job.
	State() JobState
	// SetState sets the state of the job.
	SetState(state JobState) error
	// String returns a string representation of the job.
	String() string
}
