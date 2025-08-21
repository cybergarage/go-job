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
	"fmt"
	"time"
)

func ExampleManager_scheduleJob() {
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
		WithScheduleAt(time.Now()), // Immediate scheduling is the default, so this option is redundant
		WithCompleteProcessor(func(inst Instance, res []any) {
			fmt.Printf("Result: %v\n", res)
		}),
		WithArguments(1, 2),
	)
	// Start the job manager
	mgr.Start()

	// Wait waits for all jobs to complete or terminate.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mgr.Wait(ctx)

	// Stop the job manager
	mgr.Stop()

	// Output:
	// Result: [3]
}

func ExampleManager_scheduleRegisteredJob() {
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
		WithScheduleAt(time.Now()), // Immediate scheduling is the default, so this option is redundant
		WithCompleteProcessor(func(inst Instance, res []any) {
			fmt.Printf("Result: %v\n", res)
		}),
		WithArguments(1, 2),
	)
	// Start the job manager
	mgr.Start()

	// Wait waits for all jobs to complete or terminate.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mgr.Wait(ctx)

	// Stop the job manager
	mgr.Stop()

	// Output:
	// Result: [3]
}

func ExampleManager_resizeWorkers() {
	// Create a manager with 1 worker (default)
	mgr, _ := NewManager(WithNumWorkers(1))
	mgr.Start()
	defer mgr.Stop()

	// Print the current number of workers
	fmt.Println("Initial workers:", mgr.NumWorkers())

	// Monitor queue size and scale accordingly
	query := NewQuery(
		WithQueryState(JobScheduled), // filter by scheduled state
	)
	jobs, _ := mgr.LookupInstances(query)
	queueSize := len(jobs)
	currentWorkers := mgr.NumWorkers()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if queueSize > currentWorkers*2 {
		// Scale up if queue is getting too long
		mgr.ResizeWorkers(ctx, currentWorkers+2)
	} else if queueSize == 0 && currentWorkers > 2 {
		// Scale down if no jobs queued
		mgr.ResizeWorkers(ctx, currentWorkers-1)
	}

	// Print the current number of workers
	fmt.Println("Active workers:", mgr.NumWorkers())

	// Output:
	// Initial workers: 1
	// Active workers: 1
}
