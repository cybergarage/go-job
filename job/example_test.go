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

func Example() {
	// Create a job manager
	mgr, _ := NewManager()

	// Register a job with a custom executor
	sumJob, _ := NewJob(
		WithKind("sum"),
		WithExecutor(func(a, b int) int { return a + b }),
		WithScheduleAt(time.Now()), // immediate scheduling is the default, so this option is redundant
		WithResponseHandler(func(ji Instance, res []any) {
			ji.Infof("Result: %v", res)
		}),
		WithErrorHandler(func(ji Instance, err error) error {
			ji.Errorf("Error executing job: %v", err)
			return err
		}),
	)
	mgr.RegisterJob(sumJob)

	// Schedule the registered job
	ji, _ := mgr.ScheduleRegisteredJob("sum", WithArguments(1, 2))

	// Start the job manager
	mgr.Start()

	// Wait for the job to complete
	mgr.StopWithWait()

	// Retrieve and print the job instance state history
	history, _ := mgr.ProcessHistory(ji)
	for _, record := range history {
		fmt.Println(record.State())
	}

	// Retrieve and print the job instance logs
	logs, _ := mgr.ProcessLogs(ji)
	for _, log := range logs {
		fmt.Println(log.Message())
	}

	// Output:
	// Created
	// Scheduled
	// Processing
	// Completed
	// Result: [3]
}
