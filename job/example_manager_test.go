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
	"time"
)

func ExampleManager_ScheduleJob() {
	// Create a job manager
	mgr, _ := NewManager()
	// Create a job with a custom executor
	job, _ := NewJob(
		WithKind("sum"),
		WithExecutor(func(a, b int) int { return a + b }),
	)
	// Schedule the job with the manager
	mgr.ScheduleJob(
		job,
		WithScheduleAt(time.Now()), // Schedule the job to run immediately.
		WithResponseHandler(func(inst Instance, res []any) {
			fmt.Printf("Result: %v\n", res)
		}),
		WithArguments(1, 2),
	)
	// Start the job manager
	mgr.Start()
	// Wait for the job to complete
	mgr.StopWithWait()

	// Output: Result: [3]
}

func ExampleManager_ScheduleRegisteredJob() {
	// Create a job manager
	mgr, _ := NewManager()
	// Create a job with a custom executor
	job, _ := NewJob(
		WithKind("sum"),
		WithExecutor(func(a, b int) int { return a + b }),
	)
	// Register the job with the manager
	mgr.RegisterJob(job)
	// Schedule the registered job
	mgr.ScheduleRegisteredJob(
		job.Kind(),
		WithScheduleAt(time.Now()), // Schedule the job to run immediately.
		WithResponseHandler(func(inst Instance, res []any) {
			fmt.Printf("Result: %v\n", res)
		}),
		WithArguments(1, 2),
	)
	// Start the job manager
	mgr.Start()
	// Wait for the job to complete
	mgr.StopWithWait()

	// Output: Result: [3]
}
