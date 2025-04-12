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

import (
	"context"
)

// JobManager is an interface that defines methods for managing jobs.
type JobManager interface {
	// // ProcessJob processes a job by executing it and updating its state.
	ProcessJob(ctx context.Context, job Job) error
	// // ScheduleJob schedules a job to run at a specific time or interval.
	// ScheduleJob(ctx context.Context, job Job) error
	// // CancelJob cancels a scheduled job.
	// CancelJob(ctx context.Context, job Job) error
	// 	// ListJobs lists all jobs with the specified state.
	// 	ListJobs(ctx context.Context, state JobState) ([]Job, error)
}
