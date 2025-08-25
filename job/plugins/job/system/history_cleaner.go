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

package system

import (
	"context"
	"time"

	"github.com/cybergarage/go-job/job"
	"github.com/cybergarage/go-job/job/plugins"
)

const (
	HistoryCleaner = "system.history.cleaner"
)

// NewHistoryCleaner returns a job that cleans up old job instance history records.
// The job executor accepts the following parameters:
//   - ctx: context.Context - The context for cancellation and timeout.
//   - mgr: job.Manager - The job manager to perform the cleanup operation.
//   - ji: job.Instance - The job instance representing the history cleaner job.
//   - before: time.Time - A timestamp indicating that all job instances completed before this time should be deleted.
func NewHistoryCleaner() plugins.Job {
	return plugins.NewJob(
		HistoryCleaner,
		func(ctx context.Context, mgr job.Manager, ji job.Instance, before time.Time) {
			filter := job.NewFilter(
				job.WithFilterBefore(before),
			)
			err := mgr.ClearInstanceHistory(filter)
			if err != nil {
				ji.Errorf("Failed to clear job history: %v", err)
			}
		},
	)
}
