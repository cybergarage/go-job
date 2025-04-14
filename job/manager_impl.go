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
)

type manager struct {
	store Store
}

// NewManager creates a new instance of the job manager.
func NewManager() *manager {
	return &manager{
		store: NewMemStore(),
	}
}

// ProcessJob processes a job by executing it and updating its state.
func (mgr *manager) ProcessJob(ctx context.Context, job Job) error {
	if job == nil {
		return nil
	}
	// Implement the logic to process the job
	// For example, you might want to execute the job's command or function
	// and update its state in the database or in-memory store.
	return nil
}

// ScheduleJob schedules a job to run at a specific time or interval.
func (mgr *manager) ScheduleJob(ctx context.Context, job Job) error {
	// Implement the logic to schedule the job
	// For example, you might want to add the job to a queue or a scheduler
	// that will execute it at the specified time or interval.
	return nil
}

// CancelJob cancels a scheduled job.
func (mgr *manager) CancelJob(ctx context.Context, job Job) error {
	// Implement the logic to cancel the job
	// For example, you might want to remove the job from the queue or scheduler
	// and update its state to "canceled" in the database or in-memory store.
	return nil
}

// ListJobs lists all jobs with the specified state.
func (mgr *manager) ListJobs(ctx context.Context, state JobState) ([]Job, error) {
	return mgr.store.ListJobs(ctx)
}
