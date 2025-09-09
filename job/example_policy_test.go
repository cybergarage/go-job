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

package job_test

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/cybergarage/go-job/job"
)

func ExampleWithPriority() {
	// Create a job manager
	mgr, _ := job.NewManager()

	// Create and register a job with a specific priority
	newJob, _ := job.NewJob(
		job.WithKind("sum"),
		job.WithExecutor(func(a, b int) int { return a + b }),
		job.WithPriority(0),
	)
	mgr.RegisterJob(newJob)
	fmt.Printf("Priority: %d\n", newJob.Policy().Priority())

	// Start the job manager
	mgr.Start()
	defer mgr.Stop()

	// Schedule the registered job with default job priority
	ji, _ := mgr.ScheduleRegisteredJob("sum")
	fmt.Printf("Priority: %d\n", ji.Priority())

	// Schedule the registered job with an overridden priority
	ji, _ = mgr.ScheduleRegisteredJob("sum", job.WithPriority(1))
	fmt.Printf("Priority: %d\n", ji.Priority())

	// Output:
	// Priority: 0
	// Priority: 0
	// Priority: 1
}

func ExampleWithMaxRetries() {
	// Create a job manager
	mgr, _ := job.NewManager()

	// Create and register a job with a specific max retries
	newJob, _ := job.NewJob(
		job.WithKind("sum"),
		job.WithExecutor(func(a, b int) int { return a + b }),
		job.WithMaxRetries(3),
	)
	mgr.RegisterJob(newJob)
	fmt.Printf("MaxRetries: %d\n", newJob.Policy().MaxRetries())

	// Start the job manager
	mgr.Start()
	defer mgr.Stop()

	// Schedule the registered job with default job max retries
	ji, _ := mgr.ScheduleRegisteredJob("sum")
	fmt.Printf("MaxRetries: %d\n", ji.MaxRetries())

	// Schedule the registered job with an overridden max retries
	ji, _ = mgr.ScheduleRegisteredJob("sum", job.WithMaxRetries(5))
	fmt.Printf("MaxRetries: %d\n", ji.MaxRetries())

	// Output:
	// MaxRetries: 3
	// MaxRetries: 3
	// MaxRetries: 5
}

func ExampleWithTimeout() {
	// Create a job manager
	mgr, _ := job.NewManager()

	// Create and register a job with a specific timeout
	newJob, _ := job.NewJob(
		job.WithKind("sum"),
		job.WithExecutor(func(a, b int) int { return a + b }),
		job.WithTimeout(5*time.Second),
	)
	mgr.RegisterJob(newJob)
	fmt.Printf("Timeout: %v\n", newJob.Policy().Timeout())

	// Start the job manager
	mgr.Start()
	defer mgr.Stop()

	// Schedule the registered job with default job timeout
	ji, _ := mgr.ScheduleRegisteredJob("sum")
	fmt.Printf("Timeout: %v\n", ji.Timeout())

	// Schedule the registered job with an overridden timeout
	ji, _ = mgr.ScheduleRegisteredJob("sum", job.WithTimeout(10*time.Second))
	fmt.Printf("Timeout: %v\n", ji.Timeout())

	// Output:
	// Timeout: 5s
	// Timeout: 5s
	// Timeout: 10s
}

func ExampleWithBackoffStrategy() {
	// Create and register a job with a specific backoff strategy
	job, _ := job.NewJob(
		job.WithKind("sum"),
		job.WithExecutor(func(a, b int) int { return a + b }),
		job.WithBackoffStrategy(func(ji job.Instance) time.Duration {
			// Exponential backoff
			return time.Duration(float64(ji.Attempts()) * float64(time.Second) * (0.8 + 0.4*rand.Float64()))
		}))

	fmt.Printf("BackoffStrategy: %T\n", job.Policy().BackoffStrategy())

	// Output:
	// BackoffStrategy: job.BackoffStrategy
}
