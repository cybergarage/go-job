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

func ExampleWithJitter() {
	// Create a job manager
	mgr, _ := NewManager()

	// Create and register a job with a specific jitter
	job, _ := NewJob(
		WithKind("sum"),
		WithExecutor(func(a, b int) int { return a + b }),
		WithJitter(
			func() time.Duration {
				return 100 * time.Millisecond
			},
		),
	)
	mgr.RegisterJob(job)
	fmt.Printf("Jitter: %v\n", job.Schedule().Jitter()())

	// Start the job manager
	mgr.Start()
	defer mgr.Stop()

	// Schedule the registered job with default job jitter
	ji, _ := mgr.ScheduleRegisteredJob("sum")
	fmt.Printf("Jitter: %v\n", ji.Jitter()())

	// Schedule the registered job with an overridden jitter
	ji, _ = mgr.ScheduleRegisteredJob(
		"sum",
		WithJitter(
			func() time.Duration {
				return 200 * time.Millisecond
			},
		))
	fmt.Printf("Jitter: %v\n", ji.Jitter()())

	// Output:
	// Jitter: 100ms
	// Jitter: 100ms
	// Jitter: 200ms
}
