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
	"time"

	"github.com/cybergarage/go-job/job"
)

func ExampleNewHistoryCleaner() {
	// Create a job manager
	mgr, _ := job.NewManager()

	// Schedule the job with the manager
	mgr.ScheduleRegisteredJob(
		HistoryCleaner,
		job.WithCrontabSpec("0 0 * * *"), // Every day at midnight
		job.WithArguments("?", "?", "?", time.Now().Add(-30*24*time.Hour)), // Delete history older than 30 days
	)
}

func ExampleNewLogCleaner() {
	// Create a job manager
	mgr, _ := job.NewManager()

	// Schedule the job with the manager
	mgr.ScheduleRegisteredJob(
		LogCleaner,
		job.WithCrontabSpec("0 0 * * *"), // Every day at midnight
		job.WithArguments("?", "?", "?", time.Now().Add(-30*24*time.Hour)), // Delete log older than 30 days
	)
}
