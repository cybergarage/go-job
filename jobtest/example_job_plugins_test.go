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

package jobtest

import (
	"fmt"
	"time"

	"github.com/cybergarage/go-job/job"
	"github.com/cybergarage/go-job/job/plugins/job/system"
)

func ExampleManager_scheduleRegisteredJob_systemHistoryCleaner() {
	// Create a job manager
	mgr, _ := job.NewManager()

	// Schedule the job with the manager
	_, err := mgr.ScheduleJob(
		system.NewHistoryCleaner(),
		job.WithCrontabSpec("0 0 * * *"),                                   // Every day at midnight
		job.WithArguments("?", "?", "?", time.Now().Add(-30*24*time.Hour)), // Delete history older than 30 days
	)

	fmt.Println(err)

	// Output:
	// <nil>
}

func ExampleManager_scheduleRegisteredJob_systemLogCleaner() {
	// Create a job manager
	mgr, _ := job.NewManager()

	// Schedule the job with the manager
	_, err := mgr.ScheduleJob(
		system.NewLogCleaner(),
		job.WithCrontabSpec("0 0 * * *"),                                   // Every day at midnight
		job.WithArguments("?", "?", "?", time.Now().Add(-30*24*time.Hour)), // Delete log older than 30 days
	)

	fmt.Println(err)

	// Output:
	// <nil>
}
